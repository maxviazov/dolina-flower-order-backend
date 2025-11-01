package sqlite

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/maxviazov/dolina-flower-order-backend/internal/domain"
)

type Repository struct {
	db *sql.DB
}

// Проверка соответствия интерфейсу
var _ domain.OrderRepository = (*Repository)(nil)

func NewRepository(dbPath string) (*Repository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	repo := &Repository{db: db}
	if err := repo.createTables(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *Repository) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS flowers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			mark_box TEXT NOT NULL,
			variety TEXT NOT NULL,
			length INTEGER NOT NULL,
			box_count REAL NOT NULL,
			pack_rate INTEGER NOT NULL,
			total_stems INTEGER NOT NULL,
			farm_name TEXT NOT NULL,
			truck_name TEXT NOT NULL,
			price REAL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id TEXT PRIMARY KEY,
			mark_box TEXT NOT NULL,
			customer_id TEXT NOT NULL,
			status TEXT DEFAULT 'pending',
			total_amount REAL DEFAULT 0,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			processed_at DATETIME,
			farm_order_id TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS order_items (
			id TEXT PRIMARY KEY,
			order_id TEXT NOT NULL,
			variety TEXT NOT NULL,
			length INTEGER NOT NULL,
			box_count REAL NOT NULL,
			pack_rate INTEGER NOT NULL,
			total_stems INTEGER NOT NULL,
			farm_name TEXT NOT NULL,
			truck_name TEXT NOT NULL,
			comments TEXT,
			price REAL DEFAULT 0,
			FOREIGN KEY (order_id) REFERENCES orders(id)
		)`,
	}

	for _, query := range queries {
		if _, err := r.db.Exec(query); err != nil {
			return err
		}
	}

	return r.insertTestData()
}

func (r *Repository) insertTestData() error {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM flowers").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	testFlowers := []struct {
		markBox    string
		variety    string
		length     int
		boxCount   float64
		packRate   int
		totalStems int
		farmName   string
		truckName  string
	}{
		{"VVA", "Red Naomi", 70, 10.5, 20, 210, "KENYA FARM 1", "TRUCK A"},
		{"VVA", "Freedom", 60, 8.0, 25, 200, "KENYA FARM 1", "TRUCK A"},
		{"VVA", "Explorer", 70, 12.0, 20, 240, "KENYA FARM 2", "TRUCK B"},
		{"VVA", "Avalanche", 60, 15.0, 25, 375, "KENYA FARM 2", "TRUCK B"},
		{"VVA", "Mondial", 70, 6.5, 20, 130, "KENYA FARM 3", "TRUCK C"},
		{"VVA", "Pink Floyd", 60, 9.0, 25, 225, "KENYA FARM 3", "TRUCK C"},
		{"VVA", "Rhodos", 70, 11.0, 20, 220, "KENYA FARM 1", "TRUCK A"},
		{"VVA", "Tacazzi", 60, 7.5, 25, 187, "KENYA FARM 2", "TRUCK B"},
	}

	stmt, err := r.db.Prepare(
		`
		INSERT INTO flowers (mark_box, variety, length, box_count, pack_rate, total_stems, farm_name, truck_name)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, flower := range testFlowers {
		_, err := stmt.Exec(
			flower.markBox,
			flower.variety,
			flower.length,
			flower.boxCount,
			flower.packRate,
			flower.totalStems,
			flower.farmName,
			flower.truckName,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) GetAvailableFlowers(ctx context.Context) ([]domain.Item, error) {
	rows, err := r.db.QueryContext(
		ctx, `
		SELECT mark_box, variety, length, box_count, pack_rate, total_stems, farm_name, truck_name, price
		FROM flowers
		ORDER BY variety, length
	`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flowers []domain.Item
	for rows.Next() {
		var flower domain.Item
		var markBox string
		err := rows.Scan(
			&markBox,
			&flower.Variety,
			&flower.Length,
			&flower.BoxCount,
			&flower.PackRate,
			&flower.TotalStems,
			&flower.FarmName,
			&flower.TruckName,
			&flower.Price,
		)
		if err != nil {
			return nil, err
		}
		flowers = append(flowers, flower)
	}

	return flowers, nil
}

func (r *Repository) Create(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx, `
		INSERT INTO orders (id, mark_box, customer_id, status, total_amount, notes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, order.ID, order.MarkBox, order.CustomerID, order.Status, order.TotalAmount, order.Notes, order.CreatedAt,
	)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = tx.ExecContext(
			ctx, `
			INSERT INTO order_items (id, order_id, variety, length, box_count, pack_rate, total_stems, farm_name, truck_name, comments, price)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, item.ID, order.ID, item.Variety, item.Length, item.BoxCount, item.PackRate, item.TotalStems, item.FarmName,
			item.TruckName, item.Comments, item.Price,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	var order domain.Order
	var processedAt sql.NullTime
	var farmOrderID sql.NullString

	err := r.db.QueryRowContext(
		ctx, `
		SELECT id, mark_box, customer_id, status, total_amount, notes, created_at, processed_at, farm_order_id
		FROM orders WHERE id = ?
	`, id,
	).Scan(
		&order.ID,
		&order.MarkBox,
		&order.CustomerID,
		&order.Status,
		&order.TotalAmount,
		&order.Notes,
		&order.CreatedAt,
		&processedAt,
		&farmOrderID,
	)
	if err != nil {
		return nil, err
	}

	if processedAt.Valid {
		order.ProcessedAt = &processedAt.Time
	}
	if farmOrderID.Valid {
		order.FarmOrderID = &farmOrderID.String
	}

	rows, err := r.db.QueryContext(
		ctx, `
		SELECT id, variety, length, box_count, pack_rate, total_stems, farm_name, truck_name, comments, price
		FROM order_items WHERE order_id = ?
	`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.Item
		err := rows.Scan(
			&item.ID,
			&item.Variety,
			&item.Length,
			&item.BoxCount,
			&item.PackRate,
			&item.TotalStems,
			&item.FarmName,
			&item.TruckName,
			&item.Comments,
			&item.Price,
		)
		if err != nil {
			return nil, err
		}
		item.OrderID = order.ID
		order.Items = append(order.Items, item)
	}

	return &order, nil
}

func (r *Repository) GetByStatus(ctx context.Context, status domain.OrderStatus) ([]*domain.Order, error) {
	rows, err := r.db.QueryContext(
		ctx, `
		SELECT id, mark_box, customer_id, status, total_amount, notes, created_at, processed_at, farm_order_id
		FROM orders WHERE status = ?
		ORDER BY created_at DESC
	`, status,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		var processedAt sql.NullTime
		var farmOrderID sql.NullString
		err := rows.Scan(
			&order.ID,
			&order.MarkBox,
			&order.CustomerID,
			&order.Status,
			&order.TotalAmount,
			&order.Notes,
			&order.CreatedAt,
			&processedAt,
			&farmOrderID,
		)
		if err != nil {
			return nil, err
		}
		if processedAt.Valid {
			order.ProcessedAt = &processedAt.Time
		}
		if farmOrderID.Valid {
			order.FarmOrderID = &farmOrderID.String
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func (r *Repository) Update(ctx context.Context, order *domain.Order) error {
	_, err := r.db.ExecContext(
		ctx, `
		UPDATE orders 
		SET mark_box = ?, status = ?, total_amount = ?, notes = ?, processed_at = ?, farm_order_id = ?
		WHERE id = ?
	`, order.MarkBox, order.Status, order.TotalAmount, order.Notes, order.ProcessedAt, order.FarmOrderID, order.ID,
	)
	return err
}

func (r *Repository) Close() error {
	return r.db.Close()
}
