package web

import (
	"time"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
	"github.com/RoyerHernandez/pos-system/api/internal/service"
)

// ---------------------------------------------------------------------------
// Auth
// ---------------------------------------------------------------------------

// LoginDTO represents the login request body.
type LoginDTO struct {
	Usuario  string `json:"usuario"`
	Password string `json:"password"`
}

// toDomain returns the username and password as plain values.
func (d LoginDTO) toDomain() (string, string) {
	return d.Usuario, d.Password
}

// RefreshDTO represents the token refresh request body.
type RefreshDTO struct {
	RefreshToken string `json:"refresh_token"`
}

// AuthResponseDTO is the JSON shape returned after login or token refresh.
type AuthResponseDTO struct {
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	User         UserResponseDTO `json:"user"`
}

// authResponseFromService converts a service.AuthResponse to a DTO.
func authResponseFromService(r *service.AuthResponse) AuthResponseDTO {
	return AuthResponseDTO{
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
		User:         userResponseFromDomain(r.User),
	}
}

// ---------------------------------------------------------------------------
// User
// ---------------------------------------------------------------------------

// UserResponseDTO is the JSON representation of a user (no password).
type UserResponseDTO struct {
	ID            int        `json:"id"`
	Usuario       string     `json:"usuario"`
	Nombre        string     `json:"nombre"`
	Perfil        string     `json:"perfil"`
	Foto          *string    `json:"foto"`
	UltimoLogin   *time.Time `json:"ultimo_login"`
	Estado        int        `json:"estado"`
	FechaCreacion *time.Time `json:"fecha_creacion"`
}

// userResponseFromDomain converts domain.UserResponse to a DTO.
func userResponseFromDomain(u domain.UserResponse) UserResponseDTO {
	return UserResponseDTO{
		ID:            u.ID,
		Usuario:       u.Usuario,
		Nombre:        u.Nombre,
		Perfil:        u.Perfil,
		Foto:          u.Foto,
		UltimoLogin:   u.UltimoLogin,
		Estado:        u.Estado,
		FechaCreacion: u.FechaCreacion,
	}
}

// ---------------------------------------------------------------------------
// Category
// ---------------------------------------------------------------------------

// CreateCategoryDTO represents the create category request body.
type CreateCategoryDTO struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

// UpdateCategoryDTO represents the update category request body.
type UpdateCategoryDTO struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

// CategoryResponseDTO is the JSON representation of a category.
type CategoryResponseDTO struct {
	ID            int        `json:"id"`
	Nombre        string     `json:"nombre"`
	Descripcion   *string    `json:"descripcion"`
	Estado        int        `json:"estado"`
	FechaCreacion *time.Time `json:"fecha_creacion"`
}

// categoryResponseFromDomain converts domain.Category to a DTO.
func categoryResponseFromDomain(c domain.Category) CategoryResponseDTO {
	return CategoryResponseDTO{
		ID:            c.ID,
		Nombre:        c.Nombre,
		Descripcion:   c.Descripcion,
		Estado:        c.Estado,
		FechaCreacion: c.FechaCreacion,
	}
}

// categoryListFromDomain converts a slice of domain.Category to DTOs.
func categoryListFromDomain(cats []domain.Category) []CategoryResponseDTO {
	out := make([]CategoryResponseDTO, len(cats))
	for i, c := range cats {
		out[i] = categoryResponseFromDomain(c)
	}
	return out
}

// ---------------------------------------------------------------------------
// Product
// ---------------------------------------------------------------------------

// ProductResponseDTO is the JSON representation of a product.
type ProductResponseDTO struct {
	ID            int        `json:"id"`
	Codigo        string     `json:"codigo"`
	CodigoBarras  *string    `json:"codigo_barras"`
	Descripcion   string     `json:"descripcion"`
	IDCategoria   int        `json:"id_categoria"`
	PrecioCompra  float64    `json:"precio_compra"`
	PrecioVenta   float64    `json:"precio_venta"`
	Stock         int        `json:"stock"`
	StockMinimo   int        `json:"stock_minimo"`
	Imagen        *string    `json:"imagen"`
	Estado        int        `json:"estado"`
	FechaCreacion *time.Time `json:"fecha_creacion"`
}

// productResponseFromDomain converts domain.Product to a DTO.
func productResponseFromDomain(p domain.Product) ProductResponseDTO {
	return ProductResponseDTO{
		ID:            p.ID,
		Codigo:        p.Codigo,
		CodigoBarras:  p.CodigoBarras,
		Descripcion:   p.Descripcion,
		IDCategoria:   p.IDCategoria,
		PrecioCompra:  p.PrecioCompra,
		PrecioVenta:   p.PrecioVenta,
		Stock:         p.Stock,
		StockMinimo:   p.StockMinimo,
		Imagen:        p.Imagen,
		Estado:        p.Estado,
		FechaCreacion: p.FechaCreacion,
	}
}

// productListFromDomain converts a slice of domain.Product to DTOs.
func productListFromDomain(prods []domain.Product) []ProductResponseDTO {
	out := make([]ProductResponseDTO, len(prods))
	for i, p := range prods {
		out[i] = productResponseFromDomain(p)
	}
	return out
}

// ---------------------------------------------------------------------------
// Client
// ---------------------------------------------------------------------------

// CreateClientDTO represents the create client request body.
type CreateClientDTO struct {
	Nombre          string  `json:"nombre"`
	Documento       *string `json:"documento"`
	Email           *string `json:"email"`
	Telefono        *string `json:"telefono"`
	Direccion       *string `json:"direccion"`
	FechaNacimiento *string `json:"fecha_nacimiento"`
}

// UpdateClientDTO represents the update client request body.
type UpdateClientDTO struct {
	Nombre          string  `json:"nombre"`
	Documento       *string `json:"documento"`
	Email           *string `json:"email"`
	Telefono        *string `json:"telefono"`
	Direccion       *string `json:"direccion"`
	FechaNacimiento *string `json:"fecha_nacimiento"`
}

// ClientResponseDTO is the JSON representation of a client.
type ClientResponseDTO struct {
	ID              int        `json:"id"`
	Nombre          string     `json:"nombre"`
	Documento       *string    `json:"documento"`
	Email           *string    `json:"email"`
	Telefono        *string    `json:"telefono"`
	Direccion       *string    `json:"direccion"`
	FechaNacimiento *string    `json:"fecha_nacimiento"`
	TotalCompras    float64    `json:"total_compras"`
	Estado          int        `json:"estado"`
	FechaCreacion   *time.Time `json:"fecha_creacion"`
}

// clientResponseFromDomain converts domain.Client to a DTO.
func clientResponseFromDomain(c domain.Client) ClientResponseDTO {
	return ClientResponseDTO{
		ID:              c.ID,
		Nombre:          c.Nombre,
		Documento:       c.Documento,
		Email:           c.Email,
		Telefono:        c.Telefono,
		Direccion:       c.Direccion,
		FechaNacimiento: c.FechaNacimiento,
		TotalCompras:    c.TotalCompras,
		Estado:          c.Estado,
		FechaCreacion:   c.FechaCreacion,
	}
}

// clientListFromDomain converts a slice of domain.Client to DTOs.
func clientListFromDomain(clients []domain.Client) []ClientResponseDTO {
	out := make([]ClientResponseDTO, len(clients))
	for i, c := range clients {
		out[i] = clientResponseFromDomain(c)
	}
	return out
}

// ---------------------------------------------------------------------------
// Sale
// ---------------------------------------------------------------------------

// CreateSaleDTO represents the create sale request body.
type CreateSaleDTO struct {
	IDCliente  int           `json:"id_cliente"`
	MetodoPago string        `json:"metodo_pago"`
	Productos  []SaleItemDTO `json:"productos"`
}

// SaleItemDTO represents a single product line in a sale request.
type SaleItemDTO struct {
	IDProducto int `json:"id_producto"`
	Cantidad   int `json:"cantidad"`
}

// toServiceRequest converts a CreateSaleDTO to a service.SaleRequest.
func (d CreateSaleDTO) toServiceRequest() service.SaleRequest {
	items := make([]service.SaleItemRequest, len(d.Productos))
	for i, p := range d.Productos {
		items[i] = service.SaleItemRequest{
			IDProducto: p.IDProducto,
			Cantidad:   p.Cantidad,
		}
	}
	return service.SaleRequest{
		IDCliente:  d.IDCliente,
		MetodoPago: d.MetodoPago,
		Productos:  items,
	}
}

// SaleResponseDTO is the JSON representation of a sale with details.
type SaleResponseDTO struct {
	Sale    SaleDTO       `json:"sale"`
	Details []SaleDetailDTO `json:"details"`
}

// SaleDTO is the JSON representation of a sale header.
type SaleDTO struct {
	ID          int        `json:"id"`
	IDUsuario   int        `json:"id_usuario"`
	IDCliente   *int       `json:"id_cliente"`
	CodigoVenta string     `json:"codigo_venta"`
	Subtotal    float64    `json:"subtotal"`
	Impuesto    float64    `json:"impuesto"`
	Descuento   float64    `json:"descuento"`
	Total       float64    `json:"total"`
	MetodoPago  string     `json:"metodo_pago"`
	Estado      string     `json:"estado"`
	Fecha       *time.Time `json:"fecha"`
	IDCaja      *int       `json:"id_caja"`
}

// SaleDetailDTO is the JSON representation of a sale detail line.
type SaleDetailDTO struct {
	ID             int     `json:"id"`
	IDVenta        int     `json:"id_venta"`
	IDProducto     int     `json:"id_producto"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
	Descuento      float64 `json:"descuento"`
	Subtotal       float64 `json:"subtotal"`
}

// saleFromDomain converts domain.Sale to a DTO.
func saleFromDomain(s domain.Sale) SaleDTO {
	return SaleDTO{
		ID:          s.ID,
		IDUsuario:   s.IDUsuario,
		IDCliente:   s.IDCliente,
		CodigoVenta: s.CodigoVenta,
		Subtotal:    s.Subtotal,
		Impuesto:    s.Impuesto,
		Descuento:   s.Descuento,
		Total:       s.Total,
		MetodoPago:  s.MetodoPago,
		Estado:      s.Estado,
		Fecha:       s.Fecha,
		IDCaja:      s.IDCaja,
	}
}

// saleDetailFromDomain converts domain.SaleDetail to a DTO.
func saleDetailFromDomain(d domain.SaleDetail) SaleDetailDTO {
	return SaleDetailDTO{
		ID:             d.ID,
		IDVenta:        d.IDVenta,
		IDProducto:     d.IDProducto,
		Cantidad:       d.Cantidad,
		PrecioUnitario: d.PrecioUnitario,
		Descuento:      d.Descuento,
		Subtotal:       d.Subtotal,
	}
}

// saleResponseFromService converts a service.SaleResponse to a DTO.
func saleResponseFromService(r *service.SaleResponse) SaleResponseDTO {
	details := make([]SaleDetailDTO, len(r.Details))
	for i, d := range r.Details {
		details[i] = saleDetailFromDomain(d)
	}
	return SaleResponseDTO{
		Sale:    saleFromDomain(r.Sale),
		Details: details,
	}
}

// saleListFromDomain converts a slice of domain.Sale to DTOs.
func saleListFromDomain(sales []domain.Sale) []SaleDTO {
	out := make([]SaleDTO, len(sales))
	for i, s := range sales {
		out[i] = saleFromDomain(s)
	}
	return out
}

// ---------------------------------------------------------------------------
// Inventory
// ---------------------------------------------------------------------------

// CreateMovementDTO represents the create inventory movement request body.
type CreateMovementDTO struct {
	IDProducto    int     `json:"id_producto"`
	Tipo          string  `json:"tipo"`
	Motivo        string  `json:"motivo"`
	Cantidad      int     `json:"cantidad"`
	Observaciones *string `json:"observaciones"`
}

// InventoryMovementDTO is the JSON representation of an inventory movement.
type InventoryMovementDTO struct {
	ID            int       `json:"id"`
	IDProducto    int       `json:"id_producto"`
	IDUsuario     int       `json:"id_usuario"`
	Tipo          string    `json:"tipo"`
	Motivo        string    `json:"motivo"`
	Cantidad      int       `json:"cantidad"`
	Observaciones *string   `json:"observaciones"`
	IDReferencia  *int      `json:"id_referencia"`
	Fecha         time.Time `json:"fecha"`
}

// inventoryMovementFromDomain converts domain.InventoryMovement to a DTO.
func inventoryMovementFromDomain(m domain.InventoryMovement) InventoryMovementDTO {
	return InventoryMovementDTO{
		ID:            m.ID,
		IDProducto:    m.IDProducto,
		IDUsuario:     m.IDUsuario,
		Tipo:          m.Tipo,
		Motivo:        m.Motivo,
		Cantidad:      m.Cantidad,
		Observaciones: m.Observaciones,
		IDReferencia:  m.IDReferencia,
		Fecha:         m.Fecha,
	}
}

// inventoryMovementListFromDomain converts a slice to DTOs.
func inventoryMovementListFromDomain(ms []domain.InventoryMovement) []InventoryMovementDTO {
	out := make([]InventoryMovementDTO, len(ms))
	for i, m := range ms {
		out[i] = inventoryMovementFromDomain(m)
	}
	return out
}

// ---------------------------------------------------------------------------
// Cash Register
// ---------------------------------------------------------------------------

// OpenCashRegisterDTO represents the open register request body.
type OpenCashRegisterDTO struct {
	MontoApertura float64 `json:"monto_apertura"`
}

// CashRegisterDTO is the JSON representation of a cash register.
type CashRegisterDTO struct {
	ID                 int        `json:"id"`
	IDUsuario          int        `json:"id_usuario"`
	MontoApertura      float64    `json:"monto_apertura"`
	MontoCierre        *float64   `json:"monto_cierre"`
	TotalVentas        float64    `json:"total_ventas"`
	TotalEfectivo      float64    `json:"total_efectivo"`
	TotalTarjeta       float64    `json:"total_tarjeta"`
	TotalTransferencia float64    `json:"total_transferencia"`
	Estado             string     `json:"estado"`
	FechaApertura      *time.Time `json:"fecha_apertura"`
	FechaCierre        *time.Time `json:"fecha_cierre"`
}

// cashRegisterFromDomain converts domain.CashRegister to a DTO.
func cashRegisterFromDomain(r domain.CashRegister) CashRegisterDTO {
	return CashRegisterDTO{
		ID:                 r.ID,
		IDUsuario:          r.IDUsuario,
		MontoApertura:      r.MontoApertura,
		MontoCierre:        r.MontoCierre,
		TotalVentas:        r.TotalVentas,
		TotalEfectivo:      r.TotalEfectivo,
		TotalTarjeta:       r.TotalTarjeta,
		TotalTransferencia: r.TotalTransferencia,
		Estado:             r.Estado,
		FechaApertura:      r.FechaApertura,
		FechaCierre:        r.FechaCierre,
	}
}

// cashRegisterListFromDomain converts a slice to DTOs.
func cashRegisterListFromDomain(regs []domain.CashRegister) []CashRegisterDTO {
	out := make([]CashRegisterDTO, len(regs))
	for i, r := range regs {
		out[i] = cashRegisterFromDomain(r)
	}
	return out
}

// ---------------------------------------------------------------------------
// Dashboard
// ---------------------------------------------------------------------------

// DashboardKPIsDTO is the JSON representation of dashboard KPIs.
type DashboardKPIsDTO struct {
	TodaySalesCount int                `json:"today_sales_count"`
	TodayRevenue    float64            `json:"today_revenue"`
	TotalProducts   int                `json:"total_products"`
	LowStockCount   int                `json:"low_stock_count"`
	TopProducts     []TopProductDTO    `json:"top_products"`
	RecentSales     []RecentSaleDTO    `json:"recent_sales"`
}

// TopProductDTO is the JSON representation of a top-selling product.
type TopProductDTO struct {
	ID          int     `json:"id"`
	Descripcion string  `json:"descripcion"`
	PrecioVenta float64 `json:"precio_venta"`
	Stock       int     `json:"stock"`
	TotalSold   int     `json:"total_sold"`
}

// RecentSaleDTO is the JSON representation of a recent sale.
type RecentSaleDTO struct {
	ID          int        `json:"id"`
	CodigoVenta string     `json:"codigo_venta"`
	Total       float64    `json:"total"`
	MetodoPago  string     `json:"metodo_pago"`
	Estado      string     `json:"estado"`
	Fecha       *time.Time `json:"fecha"`
}

// dashboardKPIsFromDomain converts domain.DashboardKPIs to a DTO.
func dashboardKPIsFromDomain(k *domain.DashboardKPIs) DashboardKPIsDTO {
	tops := make([]TopProductDTO, len(k.TopProducts))
	for i, p := range k.TopProducts {
		tops[i] = TopProductDTO{
			ID:          p.ID,
			Descripcion: p.Descripcion,
			PrecioVenta: p.PrecioVenta,
			Stock:       p.Stock,
			TotalSold:   p.TotalSold,
		}
	}
	recents := make([]RecentSaleDTO, len(k.RecentSales))
	for i, s := range k.RecentSales {
		recents[i] = RecentSaleDTO{
			ID:          s.ID,
			CodigoVenta: s.CodigoVenta,
			Total:       s.Total,
			MetodoPago:  s.MetodoPago,
			Estado:      s.Estado,
			Fecha:       s.Fecha,
		}
	}
	return DashboardKPIsDTO{
		TodaySalesCount: k.TodaySalesCount,
		TodayRevenue:    k.TodayRevenue,
		TotalProducts:   k.TotalProducts,
		LowStockCount:   k.LowStockCount,
		TopProducts:     tops,
		RecentSales:     recents,
	}
}

// ---------------------------------------------------------------------------
// Reports
// ---------------------------------------------------------------------------

// SalesReportDTO is the JSON representation of a sales report.
type SalesReportDTO struct {
	TotalSales    int     `json:"total_sales"`
	TotalRevenue  float64 `json:"total_revenue"`
	TotalEfectivo float64 `json:"total_efectivo"`
	TotalTarjeta  float64 `json:"total_tarjeta"`
	TotalTransfer float64 `json:"total_transferencia"`
	AvgTicket     float64 `json:"avg_ticket"`
}

// salesReportFromDomain converts domain.SalesReport to a DTO.
func salesReportFromDomain(r *domain.SalesReport) SalesReportDTO {
	return SalesReportDTO{
		TotalSales:    r.TotalSales,
		TotalRevenue:  r.TotalRevenue,
		TotalEfectivo: r.TotalEfectivo,
		TotalTarjeta:  r.TotalTarjeta,
		TotalTransfer: r.TotalTransfer,
		AvgTicket:     r.AvgTicket,
	}
}

// ProductReportDTO is the JSON representation of a product report row.
type ProductReportDTO struct {
	ID           int     `json:"id"`
	Descripcion  string  `json:"descripcion"`
	TotalSold    int     `json:"total_sold"`
	TotalRevenue float64 `json:"total_revenue"`
}

// productReportFromDomain converts domain.ProductReport to a DTO.
func productReportFromDomain(r domain.ProductReport) ProductReportDTO {
	return ProductReportDTO{
		ID:           r.ID,
		Descripcion:  r.Descripcion,
		TotalSold:    r.TotalSold,
		TotalRevenue: r.TotalRevenue,
	}
}

// productReportListFromDomain converts a slice to DTOs.
func productReportListFromDomain(rs []domain.ProductReport) []ProductReportDTO {
	out := make([]ProductReportDTO, len(rs))
	for i, r := range rs {
		out[i] = productReportFromDomain(r)
	}
	return out
}

// ClientReportDTO is the JSON representation of a client report row.
type ClientReportDTO struct {
	ID           int     `json:"id"`
	Nombre       string  `json:"nombre"`
	Documento    *string `json:"documento"`
	TotalCompras float64 `json:"total_compras"`
	NumCompras   int     `json:"num_compras"`
}

// clientReportFromDomain converts domain.ClientReport to a DTO.
func clientReportFromDomain(r domain.ClientReport) ClientReportDTO {
	return ClientReportDTO{
		ID:           r.ID,
		Nombre:       r.Nombre,
		Documento:    r.Documento,
		TotalCompras: r.TotalCompras,
		NumCompras:   r.NumCompras,
	}
}

// clientReportListFromDomain converts a slice to DTOs.
func clientReportListFromDomain(rs []domain.ClientReport) []ClientReportDTO {
	out := make([]ClientReportDTO, len(rs))
	for i, r := range rs {
		out[i] = clientReportFromDomain(r)
	}
	return out
}
