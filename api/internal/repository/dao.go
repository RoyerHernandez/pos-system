package repository

import (
	"time"

	"github.com/RoyerHernandez/pos-system/api/internal/domain"
)

// ---------------------------------------------------------------------------
// userDAO
// ---------------------------------------------------------------------------

type userDAO struct {
	ID            int        `db:"id"`
	Usuario       string     `db:"usuario"`
	Password      string     `db:"password"`
	Nombre        string     `db:"nombre"`
	Perfil        string     `db:"perfil"`
	Foto          *string    `db:"foto"`
	UltimoLogin   *time.Time `db:"ultimo_login"`
	Estado        int        `db:"estado"`
	FechaCreacion *time.Time `db:"fecha_creacion"`
}

func (d *userDAO) toDomain() domain.User {
	return domain.User{
		ID:            d.ID,
		Usuario:       d.Usuario,
		Password:      d.Password,
		Nombre:        d.Nombre,
		Perfil:        d.Perfil,
		Foto:          d.Foto,
		UltimoLogin:   d.UltimoLogin,
		Estado:        d.Estado,
		FechaCreacion: d.FechaCreacion,
	}
}

func toUserDAO(u domain.User) userDAO {
	return userDAO{
		ID:            u.ID,
		Usuario:       u.Usuario,
		Password:      u.Password,
		Nombre:        u.Nombre,
		Perfil:        u.Perfil,
		Foto:          u.Foto,
		UltimoLogin:   u.UltimoLogin,
		Estado:        u.Estado,
		FechaCreacion: u.FechaCreacion,
	}
}

// ---------------------------------------------------------------------------
// saleDAO
// ---------------------------------------------------------------------------

type saleDAO struct {
	ID          int        `db:"id"`
	IDUsuario   int        `db:"id_usuario"`
	IDCliente   *int       `db:"id_cliente"`
	CodigoVenta string     `db:"codigo_venta"`
	Subtotal    float64    `db:"subtotal"`
	Impuesto    float64    `db:"impuesto"`
	Descuento   float64    `db:"descuento"`
	Total       float64    `db:"total"`
	MetodoPago  string     `db:"metodo_pago"`
	Estado      string     `db:"estado"`
	Fecha       *time.Time `db:"fecha"`
	IDCaja      *int       `db:"id_caja"`
	IDMesa      *int       `db:"id_mesa"`
}

func (d *saleDAO) toDomain() domain.Sale {
	return domain.Sale{
		ID:          d.ID,
		IDUsuario:   d.IDUsuario,
		IDCliente:   d.IDCliente,
		CodigoVenta: d.CodigoVenta,
		Subtotal:    d.Subtotal,
		Impuesto:    d.Impuesto,
		Descuento:   d.Descuento,
		Total:       d.Total,
		MetodoPago:  d.MetodoPago,
		Estado:      d.Estado,
		Fecha:       d.Fecha,
		IDCaja:      d.IDCaja,
		IDMesa:      d.IDMesa,
	}
}

func toSaleDAO(s domain.Sale) saleDAO {
	return saleDAO{
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
		IDMesa:      s.IDMesa,
	}
}

// ---------------------------------------------------------------------------
// tableDAO
// ---------------------------------------------------------------------------

type tableDAO struct {
	ID                 int        `db:"id"`
	Numero             int        `db:"numero"`
	Nombre             *string    `db:"nombre"`
	Capacidad          int        `db:"capacidad"`
	Estado             string     `db:"estado"`
	IDVentaActiva      *int       `db:"id_venta_activa"`
	IDMesero           *int       `db:"id_mesero"`
	FechaCreacion      *time.Time `db:"fecha_creacion"`
	FechaActualizacion *time.Time `db:"fecha_actualizacion"`
}

func (d *tableDAO) toDomain() domain.Table {
	return domain.Table{
		ID:                 d.ID,
		Numero:             d.Numero,
		Nombre:             d.Nombre,
		Capacidad:          d.Capacidad,
		Estado:             d.Estado,
		IDVentaActiva:      d.IDVentaActiva,
		IDMesero:           d.IDMesero,
		FechaCreacion:      d.FechaCreacion,
		FechaActualizacion: d.FechaActualizacion,
	}
}

func toTableDAO(t domain.Table) tableDAO {
	return tableDAO{
		ID:                 t.ID,
		Numero:             t.Numero,
		Nombre:             t.Nombre,
		Capacidad:          t.Capacidad,
		Estado:             t.Estado,
		IDVentaActiva:      t.IDVentaActiva,
		IDMesero:           t.IDMesero,
		FechaCreacion:      t.FechaCreacion,
		FechaActualizacion: t.FechaActualizacion,
	}
}

// ---------------------------------------------------------------------------
// saleDetailDAO
// ---------------------------------------------------------------------------

type saleDetailDAO struct {
	ID             int     `db:"id"`
	IDVenta        int     `db:"id_venta"`
	IDProducto     int     `db:"id_producto"`
	Cantidad       int     `db:"cantidad"`
	PrecioUnitario float64 `db:"precio_unitario"`
	Descuento      float64 `db:"descuento"`
	Subtotal       float64 `db:"subtotal"`
}

func (d *saleDetailDAO) toDomain() domain.SaleDetail {
	return domain.SaleDetail{
		ID:             d.ID,
		IDVenta:        d.IDVenta,
		IDProducto:     d.IDProducto,
		Cantidad:       d.Cantidad,
		PrecioUnitario: d.PrecioUnitario,
		Descuento:      d.Descuento,
		Subtotal:       d.Subtotal,
	}
}

func toSaleDetailDAO(s domain.SaleDetail) saleDetailDAO {
	return saleDetailDAO{
		ID:             s.ID,
		IDVenta:        s.IDVenta,
		IDProducto:     s.IDProducto,
		Cantidad:       s.Cantidad,
		PrecioUnitario: s.PrecioUnitario,
		Descuento:      s.Descuento,
		Subtotal:       s.Subtotal,
	}
}

// ---------------------------------------------------------------------------
// productDAO
// ---------------------------------------------------------------------------

type productDAO struct {
	ID            int        `db:"id"`
	Codigo        string     `db:"codigo"`
	CodigoBarras  *string    `db:"codigo_barras"`
	Descripcion   string     `db:"descripcion"`
	IDCategoria   int        `db:"id_categoria"`
	PrecioCompra  float64    `db:"precio_compra"`
	PrecioVenta   float64    `db:"precio_venta"`
	Stock         int        `db:"stock"`
	StockMinimo   int        `db:"stock_minimo"`
	Imagen        *string    `db:"imagen"`
	Estado        int        `db:"estado"`
	FechaCreacion *time.Time `db:"fecha_creacion"`
}

func (d *productDAO) toDomain() domain.Product {
	return domain.Product{
		ID:            d.ID,
		Codigo:        d.Codigo,
		CodigoBarras:  d.CodigoBarras,
		Descripcion:   d.Descripcion,
		IDCategoria:   d.IDCategoria,
		PrecioCompra:  d.PrecioCompra,
		PrecioVenta:   d.PrecioVenta,
		Stock:         d.Stock,
		StockMinimo:   d.StockMinimo,
		Imagen:        d.Imagen,
		Estado:        d.Estado,
		FechaCreacion: d.FechaCreacion,
	}
}

func toProductDAO(p domain.Product) productDAO {
	return productDAO{
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

// ---------------------------------------------------------------------------
// clientDAO
// ---------------------------------------------------------------------------

type clientDAO struct {
	ID              int        `db:"id"`
	Nombre          string     `db:"nombre"`
	Documento       *string    `db:"documento"`
	Email           *string    `db:"email"`
	Telefono        *string    `db:"telefono"`
	Direccion       *string    `db:"direccion"`
	FechaNacimiento *string    `db:"fecha_nacimiento"`
	TotalCompras    float64    `db:"total_compras"`
	Estado          int        `db:"estado"`
	FechaCreacion   *time.Time `db:"fecha_creacion"`
}

func (d *clientDAO) toDomain() domain.Client {
	return domain.Client{
		ID:              d.ID,
		Nombre:          d.Nombre,
		Documento:       d.Documento,
		Email:           d.Email,
		Telefono:        d.Telefono,
		Direccion:       d.Direccion,
		FechaNacimiento: d.FechaNacimiento,
		TotalCompras:    d.TotalCompras,
		Estado:          d.Estado,
		FechaCreacion:   d.FechaCreacion,
	}
}

func toClientDAO(c domain.Client) clientDAO {
	return clientDAO{
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

// ---------------------------------------------------------------------------
// categoryDAO
// ---------------------------------------------------------------------------

type categoryDAO struct {
	ID            int        `db:"id"`
	Nombre        string     `db:"nombre"`
	Descripcion   *string    `db:"descripcion"`
	Estado        int        `db:"estado"`
	FechaCreacion *time.Time `db:"fecha_creacion"`
}

func (d *categoryDAO) toDomain() domain.Category {
	return domain.Category{
		ID:            d.ID,
		Nombre:        d.Nombre,
		Descripcion:   d.Descripcion,
		Estado:        d.Estado,
		FechaCreacion: d.FechaCreacion,
	}
}

func toCategoryDAO(c domain.Category) categoryDAO {
	return categoryDAO{
		ID:            c.ID,
		Nombre:        c.Nombre,
		Descripcion:   c.Descripcion,
		Estado:        c.Estado,
		FechaCreacion: c.FechaCreacion,
	}
}

// ---------------------------------------------------------------------------
// cashRegisterDAO
// ---------------------------------------------------------------------------

type cashRegisterDAO struct {
	ID                 int        `db:"id"`
	IDUsuario          int        `db:"id_usuario"`
	MontoApertura      float64    `db:"monto_apertura"`
	MontoCierre        *float64   `db:"monto_cierre"`
	TotalVentas        float64    `db:"total_ventas"`
	TotalEfectivo      float64    `db:"total_efectivo"`
	TotalTarjeta       float64    `db:"total_tarjeta"`
	TotalTransferencia float64    `db:"total_transferencia"`
	Estado             string     `db:"estado"`
	FechaApertura      *time.Time `db:"fecha_apertura"`
	FechaCierre        *time.Time `db:"fecha_cierre"`
}

func (d *cashRegisterDAO) toDomain() domain.CashRegister {
	return domain.CashRegister{
		ID:                 d.ID,
		IDUsuario:          d.IDUsuario,
		MontoApertura:      d.MontoApertura,
		MontoCierre:        d.MontoCierre,
		TotalVentas:        d.TotalVentas,
		TotalEfectivo:      d.TotalEfectivo,
		TotalTarjeta:       d.TotalTarjeta,
		TotalTransferencia: d.TotalTransferencia,
		Estado:             d.Estado,
		FechaApertura:      d.FechaApertura,
		FechaCierre:        d.FechaCierre,
	}
}

func toCashRegisterDAO(cr domain.CashRegister) cashRegisterDAO {
	return cashRegisterDAO{
		ID:                 cr.ID,
		IDUsuario:          cr.IDUsuario,
		MontoApertura:      cr.MontoApertura,
		MontoCierre:        cr.MontoCierre,
		TotalVentas:        cr.TotalVentas,
		TotalEfectivo:      cr.TotalEfectivo,
		TotalTarjeta:       cr.TotalTarjeta,
		TotalTransferencia: cr.TotalTransferencia,
		Estado:             cr.Estado,
		FechaApertura:      cr.FechaApertura,
		FechaCierre:        cr.FechaCierre,
	}
}

// ---------------------------------------------------------------------------
// inventoryMovementDAO
// ---------------------------------------------------------------------------

type inventoryMovementDAO struct {
	ID            int       `db:"id"`
	IDProducto    int       `db:"id_producto"`
	IDUsuario     int       `db:"id_usuario"`
	Tipo          string    `db:"tipo"`
	Motivo        string    `db:"motivo"`
	Cantidad      int       `db:"cantidad"`
	Observaciones *string   `db:"observaciones"`
	IDReferencia  *int      `db:"id_referencia"`
	Fecha         time.Time `db:"fecha"`
}

func (d *inventoryMovementDAO) toDomain() domain.InventoryMovement {
	return domain.InventoryMovement{
		ID:            d.ID,
		IDProducto:    d.IDProducto,
		IDUsuario:     d.IDUsuario,
		Tipo:          d.Tipo,
		Motivo:        d.Motivo,
		Cantidad:      d.Cantidad,
		Observaciones: d.Observaciones,
		IDReferencia:  d.IDReferencia,
		Fecha:         d.Fecha,
	}
}

func toInventoryMovementDAO(m domain.InventoryMovement) inventoryMovementDAO {
	return inventoryMovementDAO{
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

// ---------------------------------------------------------------------------
// topProductDAO (dashboard)
// ---------------------------------------------------------------------------

type topProductDAO struct {
	ID          int     `db:"id"`
	Descripcion string  `db:"descripcion"`
	PrecioVenta float64 `db:"precio_venta"`
	Stock       int     `db:"stock"`
	TotalSold   int     `db:"total_sold"`
}

func (d *topProductDAO) toDomain() domain.TopProduct {
	return domain.TopProduct{
		ID:          d.ID,
		Descripcion: d.Descripcion,
		PrecioVenta: d.PrecioVenta,
		Stock:       d.Stock,
		TotalSold:   d.TotalSold,
	}
}

// ---------------------------------------------------------------------------
// recentSaleDAO (dashboard)
// ---------------------------------------------------------------------------

type recentSaleDAO struct {
	ID          int        `db:"id"`
	CodigoVenta string     `db:"codigo_venta"`
	Total       float64    `db:"total"`
	MetodoPago  string     `db:"metodo_pago"`
	Estado      string     `db:"estado"`
	Fecha       *time.Time `db:"fecha"`
}

func (d *recentSaleDAO) toDomain() domain.RecentSale {
	return domain.RecentSale{
		ID:          d.ID,
		CodigoVenta: d.CodigoVenta,
		Total:       d.Total,
		MetodoPago:  d.MetodoPago,
		Estado:      d.Estado,
		Fecha:       d.Fecha,
	}
}

// ---------------------------------------------------------------------------
// saleReportDAO (reports)
// ---------------------------------------------------------------------------

type saleReportDAO struct {
	TotalSales    int     `db:"total_sales"`
	TotalRevenue  float64 `db:"total_revenue"`
	TotalEfectivo float64 `db:"total_efectivo"`
	TotalTarjeta  float64 `db:"total_tarjeta"`
	TotalTransfer float64 `db:"total_transferencia"`
	AvgTicket     float64 `db:"avg_ticket"`
}

func (d *saleReportDAO) toDomain() domain.SalesReport {
	return domain.SalesReport{
		TotalSales:    d.TotalSales,
		TotalRevenue:  d.TotalRevenue,
		TotalEfectivo: d.TotalEfectivo,
		TotalTarjeta:  d.TotalTarjeta,
		TotalTransfer: d.TotalTransfer,
		AvgTicket:     d.AvgTicket,
	}
}

// ---------------------------------------------------------------------------
// productReportDAO (reports)
// ---------------------------------------------------------------------------

type productReportDAO struct {
	ID           int     `db:"id"`
	Descripcion  string  `db:"descripcion"`
	TotalSold    int     `db:"total_sold"`
	TotalRevenue float64 `db:"total_revenue"`
}

func (d *productReportDAO) toDomain() domain.ProductReport {
	return domain.ProductReport{
		ID:           d.ID,
		Descripcion:  d.Descripcion,
		TotalSold:    d.TotalSold,
		TotalRevenue: d.TotalRevenue,
	}
}

// ---------------------------------------------------------------------------
// clientReportDAO (reports)
// ---------------------------------------------------------------------------

type clientReportDAO struct {
	ID           int     `db:"id"`
	Nombre       string  `db:"nombre"`
	Documento    *string `db:"documento"`
	TotalCompras float64 `db:"total_compras"`
	NumCompras   int     `db:"num_compras"`
}

func (d *clientReportDAO) toDomain() domain.ClientReport {
	return domain.ClientReport{
		ID:           d.ID,
		Nombre:       d.Nombre,
		Documento:    d.Documento,
		TotalCompras: d.TotalCompras,
		NumCompras:   d.NumCompras,
	}
}
