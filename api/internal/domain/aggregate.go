package domain

import "time"

type User struct {
	ID            int
	Usuario       string
	Password      string
	Nombre        string
	Perfil        string
	Foto          *string
	UltimoLogin   *time.Time
	Estado        int
	FechaCreacion *time.Time
}

type Sale struct {
	ID          int
	IDUsuario   int
	IDCliente   *int
	CodigoVenta string
	Subtotal    float64
	Impuesto    float64
	Descuento   float64
	Total       float64
	MetodoPago  string
	Estado      string
	Fecha       *time.Time
	IDCaja      *int
	IDMesa      *int
}

type Table struct {
	ID                 int
	Numero             int
	Nombre             *string
	Capacidad          int
	Estado             string
	IDVentaActiva      *int
	IDMesero           *int
	FechaCreacion      *time.Time
	FechaActualizacion *time.Time
}

type SaleDetail struct {
	ID             int
	IDVenta        int
	IDProducto     int
	Cantidad       int
	PrecioUnitario float64
	Descuento      float64
	Subtotal       float64
}

type Product struct {
	ID            int
	Codigo        string
	CodigoBarras  *string
	Descripcion   string
	IDCategoria   int
	PrecioCompra  float64
	PrecioVenta   float64
	Stock         int
	StockMinimo   int
	Imagen        *string
	Estado        int
	FechaCreacion *time.Time
}

type Client struct {
	ID              int
	Nombre          string
	Documento       *string
	Email           *string
	Telefono        *string
	Direccion       *string
	FechaNacimiento *string
	TotalCompras    float64
	Estado          int
	FechaCreacion   *time.Time
}

type Category struct {
	ID            int
	Nombre        string
	Descripcion   *string
	Estado        int
	FechaCreacion *time.Time
}

type CashRegister struct {
	ID                 int
	IDUsuario          int
	MontoApertura      float64
	MontoCierre        *float64
	TotalVentas        float64
	TotalEfectivo      float64
	TotalTarjeta       float64
	TotalTransferencia float64
	Estado             string
	FechaApertura      *time.Time
	FechaCierre        *time.Time
}

type InventoryMovement struct {
	ID            int
	IDProducto    int
	IDUsuario     int
	Tipo          string
	Motivo        string
	Cantidad      int
	Observaciones *string
	IDReferencia  *int
	Fecha         time.Time
}

type DashboardKPIs struct {
	TodaySalesCount int
	TodayRevenue    float64
	TotalProducts   int
	LowStockCount   int
	TopProducts     []TopProduct
	RecentSales     []RecentSale
}

type TopProduct struct {
	ID          int
	Descripcion string
	PrecioVenta float64
	Stock       int
	TotalSold   int
}

type RecentSale struct {
	ID          int
	CodigoVenta string
	Total       float64
	MetodoPago  string
	Estado      string
	Fecha       *time.Time
}

type SalesReport struct {
	TotalSales    int
	TotalRevenue  float64
	TotalEfectivo float64
	TotalTarjeta  float64
	TotalTransfer float64
	AvgTicket     float64
}

type ProductReport struct {
	ID           int
	Descripcion  string
	TotalSold    int
	TotalRevenue float64
}

type ClientReport struct {
	ID           int
	Nombre       string
	Documento    *string
	TotalCompras float64
	NumCompras   int
}
