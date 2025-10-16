package http

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"time"

	"github.com/aronipurwanto/go-download-csv/internal/domain/transaction"
	"github.com/aronipurwanto/go-download-csv/internal/middleware"
	"github.com/aronipurwanto/go-download-csv/internal/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// ---- config & keys

const (
	createLocalKey = "create_tx_body"
	updateLocalKey = "update_tx_body"

	defaultTimeout = 5 * time.Second
	defaultPage    = 1
	defaultSize    = 10
	maxSize        = 100
)

// ---- controller

type TransactionController struct {
	svc     transaction.Service
	timeout time.Duration
}

func NewTransactionController(svc transaction.Service) *TransactionController {
	return &TransactionController{svc: svc, timeout: defaultTimeout}
}

func (h *TransactionController) withCtx(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Context(), h.timeout)
}

// RegisterTransactionRoutes keeps backward-compatible signature
func RegisterTransactionRoutes(r fiber.Router, svc transaction.Service) {
	NewTransactionController(svc).Register(r)
}

func (h *TransactionController) Register(r fiber.Router) {
	g := r.Group("/transactions")

	// POST /v1/transactions
	g.Post("/",
		middleware.ValidateBody[transaction.CreateRequest](transaction.ValidateCreate, createLocalKey),
		h.create,
	)

	// GET /v1/transactions/:id
	g.Get("/:id", h.getByID)

	// GET /v1/transactions?page=1&size=10
	g.Get("/", h.list)

	// PUT /v1/transactions/:id
	g.Put("/:id",
		middleware.ValidateBody[transaction.UpdateRequest]((transaction.UpdateRequest).Validate, updateLocalKey),
		h.update,
	)

	// DELETE /v1/transactions/:id
	g.Delete("/:id", h.delete)

	g.Get("/export.csv", h.export) // <-- NEW
}

// ---- handlers

func (h *TransactionController) create(c *fiber.Ctx) error {
	req := c.Locals(createLocalKey).(transaction.CreateRequest)
	ctx, cancel := h.withCtx(c)
	defer cancel()

	res, err := h.svc.Create(ctx, req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.Created(c, res)
}

func (h *TransactionController) getByID(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := h.withCtx(c)
	defer cancel()

	res, err := h.svc.Get(ctx, id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}
	return response.Success(c, res, nil)
}

func (h *TransactionController) list(c *fiber.Ctx) error {
	page, size := parsePagination(c)
	ctx, cancel := h.withCtx(c)
	defer cancel()

	items, pageN, total, err := h.svc.List(ctx, page, size)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	meta := fiber.Map{"page": pageN, "size": size, "total": total}
	return response.Success(c, items, meta)
}

func (h *TransactionController) update(c *fiber.Ctx) error {
	id := c.Params("id")
	req := c.Locals(updateLocalKey).(transaction.UpdateRequest)

	ctx, cancel := h.withCtx(c)
	defer cancel()

	res, err := h.svc.Update(ctx, id, req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	return response.Success(c, res, nil)
}

func (h *TransactionController) delete(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := h.withCtx(c)
	defer cancel()

	if err := h.svc.Delete(ctx, id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.Map{"deleted": id}, nil)
}

// ---- helpers

func parsePagination(c *fiber.Ctx) (int, int) {
	page, _ := strconv.Atoi(c.Query("page", strconv.Itoa(defaultPage)))
	size, _ := strconv.Atoi(c.Query("size", strconv.Itoa(defaultSize)))

	if page < 1 {
		page = defaultPage
	}
	if size < 1 {
		size = defaultSize
	}
	if size > maxSize {
		size = maxSize
	}
	return page, size
}

func (h *TransactionController) export(c *fiber.Ctx) error {
	// --- parse filters from/to (YYYY-MM-DD atau RFC3339)
	parseDate := func(s string) (time.Time, error) {
		if s == "" {
			return time.Time{}, nil
		}
		if len(s) == 10 {
			return time.Parse("2006-01-02", s)
		}
		return time.Parse(time.RFC3339, s)
	}
	from, _ := parseDate(c.Query("from", ""))
	to, _ := parseDate(c.Query("to", ""))

	// --- ambil semua data via pagination (tetap pakai Service.List)
	ctx, cancel := h.withCtx(c)
	defer cancel()

	page, size := 1, 500
	var all []transaction.Response
	for {
		items, _, total, err := h.svc.List(ctx, page, size)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error())
		}
		for _, it := range items {
			if (!from.IsZero() && it.TransactionDate.Before(from)) ||
				(!to.IsZero() && it.TransactionDate.After(to)) {
				continue
			}
			all = append(all, it)
		}
		if int64(page*size) >= total || len(items) == 0 {
			break
		}
		page++
	}

	// --- kalkulasi ukuran & jumlah part dengan memperhitungkan header per part
	const chunkLimit = 10 * 1024 // 10KB
	totalBytes := estimateCSVBytes(all)
	headerBytes := estimateCSVHeaderBytes()
	// tiap part akan memiliki header sendiri, jadi kira numParts dengan overhead header
	numParts := int(math.Ceil((float64(totalBytes) + float64(headerBytes)) / (float64(chunkLimit) + float64(headerBytes))))
	if numParts < 1 {
		numParts = 1
	}

	// --- jika user meminta bagian tertentu (?part=N) => stream CSV bagian N
	if pStr := c.Query("part", ""); pStr != "" {
		part, err := strconv.Atoi(pStr)
		if err != nil || part < 1 || part > numParts {
			return response.Error(c, fiber.StatusBadRequest, "invalid part")
		}

		start, end := splitRange(len(all), numParts, part)
		fname := fmt.Sprintf("transactions_part_%d_of_%d.csv", part, numParts)
		if !from.IsZero() || !to.IsZero() {
			f := from.Format("2006-01-02")
			t := to.Format("2006-01-02")
			if f == "0001-01-01" {
				f = "all"
			}
			if t == "0001-01-01" {
				t = "all"
			}
			fname = fmt.Sprintf("transactions_%s_to_%s_part_%d_of_%d.csv", f, t, part, numParts)
		}

		// header download
		c.Type("csv")                      // Content-Type: text/csv
		c.Set("Cache-Control", "no-store") // jangan cache
		c.Attachment(fname)

		// optional: BOM untuk Excel Windows
		if c.Query("excel") == "true" {
			_, _ = c.Write([]byte{0xEF, 0xBB, 0xBF})
		}

		w := csv.NewWriter(c)
		defer w.Flush()

		if err := writeCSVHeader(w); err != nil {
			return err
		}
		for _, it := range all[start:end] {
			if err := writeCSVRow(w, it); err != nil {
				return err
			}
		}
		return nil
	}

	// --- mode MANIFEST atau SINGLE
	if numParts == 1 && totalBytes <= chunkLimit {
		// kirim 1 file
		fname := "transactions.csv"
		if !from.IsZero() || !to.IsZero() {
			f := from.Format("2006-01-02")
			t := to.Format("2006-01-02")
			if f == "0001-01-01" {
				f = "all"
			}
			if t == "0001-01-01" {
				t = "all"
			}
			fname = "transactions_" + f + "_to_" + t + ".csv"
		}
		c.Type("csv")
		c.Set("Cache-Control", "no-store")
		c.Attachment(fname)

		if c.Query("excel") == "true" {
			_, _ = c.Write([]byte{0xEF, 0xBB, 0xBF})
		}

		w := csv.NewWriter(c)
		defer w.Flush()

		if err := writeCSVHeader(w); err != nil {
			return err
		}
		for _, it := range all {
			if err := writeCSVRow(w, it); err != nil {
				return err
			}
		}
		return nil
	}

	// besar dari 10KB => bagi jadi beberapa link (manifest JSON)
	base := c.BaseURL() + c.Path()
	links := make([]string, 0, numParts)
	for i := 1; i <= numParts; i++ {
		q := url.Values{}
		if !from.IsZero() {
			q.Set("from", from.Format("2006-01-02"))
		}
		if !to.IsZero() {
			q.Set("to", to.Format("2006-01-02"))
		}
		q.Set("part", strconv.Itoa(i))
		// bawa flag excel jika ada
		if c.Query("excel") == "true" {
			q.Set("excel", "true")
		}
		links = append(links, base+"?"+q.Encode())
	}

	meta := fiber.Map{
		"total_bytes_estimate": totalBytes,
		"chunk_limit_bytes":    chunkLimit,
		"num_parts":            numParts,
	}
	return response.Success(c, fiber.Map{"links": links}, meta)
}

// hitung ukuran CSV dengan menulis ke buffer (estimasi akurat)
func estimateCSVBytes(rows []transaction.Response) int {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = writeCSVHeader(w)
	for _, it := range rows {
		_ = writeCSVRow(w, it)
	}
	w.Flush()
	return buf.Len()
}

func writeCSVHeader(w *csv.Writer) error {
	head := []string{
		"Transaction ID",
		"No Ref",
		"Order Type Code",
		"Order Type Name",
		"Transaction Type Code",
		"Transaction Type Name",
		"Transaction Date",
		"From Account Number",
		"From Account Name",
		"From Account Product Name",
		"To Account Number",
		"To Account Name",
		"To Account Product Name",
		"Amount",
		"Status",
		"Description",
		"Method",
		"Currency",
		"Metadata",
	}
	return w.Write(head)
}

func writeCSVRow(w *csv.Writer, it transaction.Response) error {
	row := []string{
		it.TransactionID,
		it.NoRef,
		it.OrderTypeCode,
		it.OrderTypeName,
		it.TransactionTypeCode,
		it.TransactionTypeName,
		it.TransactionDate.Format(time.RFC3339),
		it.FromAccountNumber,
		it.FromAccountName,
		it.FromAccountProductName,
		it.ToAccountNumber,
		it.ToAccountName,
		it.ToAccountProductName,
		strconv.FormatFloat(it.Amount, 'f', -1, 64),
		it.Status,
		it.Description,
		it.Method,
		it.Currency,
		string(it.Metadata),
	}
	return w.Write(row)
}

// bagi range data menjadi numParts bagian yang kira-kira sama rata
func splitRange(totalRows, numParts, part int) (start, end int) {
	if numParts <= 1 || totalRows == 0 {
		return 0, totalRows
	}
	// floor division per bagian, sisa dibagi ke bagian awal
	base := totalRows / numParts
	rem := totalRows % numParts
	start = (part - 1) * base
	start += min(rem, part-1) // distribusi sisa
	end = start + base
	if part <= rem {
		end++
	}
	if end > totalRows {
		end = totalRows
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func estimateCSVHeaderBytes() int {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = writeCSVHeader(w)
	w.Flush()
	return buf.Len()
}
