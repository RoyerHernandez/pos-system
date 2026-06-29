<?php

require_once "connection.php";

class ReportModel{

	/*=============================================
	SALES BY DATE RANGE (detailed)
	=============================================*/

	static public function mdlSalesByDateRange($startDate, $endDate){

		$stmt = Connection::connect()->prepare("SELECT v.id, v.codigo_venta, v.fecha, u.nombre as vendedor, c.nombre as cliente, v.metodo_pago, v.subtotal, v.impuesto, v.descuento, v.total, v.estado FROM ventas v LEFT JOIN usuarios u ON v.id_usuario = u.id LEFT JOIN clientes c ON v.id_cliente = c.id WHERE DATE(v.fecha) BETWEEN :startDate AND :endDate ORDER BY v.fecha DESC");

		$stmt -> bindParam(":startDate", $startDate, PDO::PARAM_STR);
		$stmt -> bindParam(":endDate", $endDate, PDO::PARAM_STR);

		$stmt -> execute();

		return $stmt -> fetchAll();

	}

	/*=============================================
	SALES SUMMARY BY DATE RANGE
	=============================================*/

	static public function mdlSalesSummary($startDate, $endDate){

		$stmt = Connection::connect()->prepare("SELECT COUNT(*) as total_count, COALESCE(SUM(total), 0) as total_amount, COALESCE(SUM(CASE WHEN estado = 'completada' THEN total ELSE 0 END), 0) as completed_amount, SUM(CASE WHEN estado = 'completada' THEN 1 ELSE 0 END) as completed_count, SUM(CASE WHEN estado = 'cancelada' THEN 1 ELSE 0 END) as cancelled_count FROM ventas WHERE DATE(fecha) BETWEEN :startDate AND :endDate");

		$stmt -> bindParam(":startDate", $startDate, PDO::PARAM_STR);
		$stmt -> bindParam(":endDate", $endDate, PDO::PARAM_STR);

		$stmt -> execute();

		return $stmt -> fetch();

	}

	/*=============================================
	TOP PRODUCTS BY DATE RANGE
	=============================================*/

	static public function mdlTopProducts($startDate, $endDate, $limit){

		$stmt = Connection::connect()->prepare("SELECT p.codigo, p.descripcion, SUM(d.cantidad) as total_qty, SUM(d.subtotal) as total_revenue, AVG(d.precio_unitario) as avg_price FROM detalle_ventas d INNER JOIN ventas v ON d.id_venta = v.id INNER JOIN productos p ON d.id_producto = p.id WHERE DATE(v.fecha) BETWEEN :startDate AND :endDate AND v.estado = 'completada' GROUP BY d.id_producto ORDER BY total_qty DESC LIMIT $limit");

		$stmt -> bindParam(":startDate", $startDate, PDO::PARAM_STR);
		$stmt -> bindParam(":endDate", $endDate, PDO::PARAM_STR);

		$stmt -> execute();

		return $stmt -> fetchAll();

	}

	/*=============================================
	SALES BY PAYMENT METHOD (date range)
	=============================================*/

	static public function mdlSalesByPaymentMethod($startDate, $endDate){

		$stmt = Connection::connect()->prepare("SELECT metodo_pago, COUNT(*) as total_count, COALESCE(SUM(total), 0) as total_amount FROM ventas WHERE DATE(fecha) BETWEEN :startDate AND :endDate AND estado = 'completada' GROUP BY metodo_pago ORDER BY total_amount DESC");

		$stmt -> bindParam(":startDate", $startDate, PDO::PARAM_STR);
		$stmt -> bindParam(":endDate", $endDate, PDO::PARAM_STR);

		$stmt -> execute();

		return $stmt -> fetchAll();

	}

	/*=============================================
	CASH REGISTER REPORT (closed registers in date range)
	=============================================*/

	static public function mdlCashRegisterReport($startDate, $endDate){

		$stmt = Connection::connect()->prepare("SELECT c.id, u.nombre as usuario, c.monto_apertura, c.monto_cierre, c.total_ventas, c.total_efectivo, c.total_tarjeta, c.total_transferencia, c.estado, c.fecha_apertura, c.fecha_cierre FROM caja c LEFT JOIN usuarios u ON c.id_usuario = u.id WHERE DATE(c.fecha_apertura) BETWEEN :startDate AND :endDate ORDER BY c.fecha_apertura DESC");

		$stmt -> bindParam(":startDate", $startDate, PDO::PARAM_STR);
		$stmt -> bindParam(":endDate", $endDate, PDO::PARAM_STR);

		$stmt -> execute();

		return $stmt -> fetchAll();

	}

	/*=============================================
	DAILY SALES BREAKDOWN (for chart in reports)
	=============================================*/

	static public function mdlDailySalesBreakdown($startDate, $endDate){

		$stmt = Connection::connect()->prepare("SELECT DATE(fecha) as sale_date, COUNT(*) as total_count, COALESCE(SUM(total), 0) as total_amount FROM ventas WHERE DATE(fecha) BETWEEN :startDate AND :endDate AND estado = 'completada' GROUP BY DATE(fecha) ORDER BY sale_date ASC");

		$stmt -> bindParam(":startDate", $startDate, PDO::PARAM_STR);
		$stmt -> bindParam(":endDate", $endDate, PDO::PARAM_STR);

		$stmt -> execute();

		return $stmt -> fetchAll();

	}

}
