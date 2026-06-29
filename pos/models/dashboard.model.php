<?php

require_once "connection.php";

class DashboardModel{

	/*=============================================
	TODAY'S SALES (count + total)
	=============================================*/

	static public function mdlTodaySales(){

		$stmt = Connection::connect()->prepare("SELECT COUNT(*) as total_count, COALESCE(SUM(total), 0) as total_amount FROM ventas WHERE DATE(fecha) = CURDATE() AND estado = 'completada'");

		$stmt -> execute();

		return $stmt -> fetch();

	}

	/*=============================================
	LOW STOCK PRODUCTS
	=============================================*/

	static public function mdlLowStockProducts(){

		$stmt = Connection::connect()->prepare("SELECT COUNT(*) as total FROM productos WHERE stock <= stock_minimo AND estado = 1");

		$stmt -> execute();

		return $stmt -> fetch();

	}

	/*=============================================
	LOW STOCK PRODUCTS LIST
	=============================================*/

	static public function mdlLowStockProductsList(){

		$stmt = Connection::connect()->prepare("SELECT codigo, descripcion, stock, stock_minimo FROM productos WHERE stock <= stock_minimo AND estado = 1 ORDER BY stock ASC LIMIT 10");

		$stmt -> execute();

		return $stmt -> fetchAll();

	}

	/*=============================================
	TOTAL ACTIVE CLIENTS
	=============================================*/

	static public function mdlTotalClients(){

		$stmt = Connection::connect()->prepare("SELECT COUNT(*) as total FROM clientes WHERE estado = 1");

		$stmt -> execute();

		return $stmt -> fetch();

	}

	/*=============================================
	SALES LAST 7 DAYS (for bar chart)
	=============================================*/

	static public function mdlSalesLast7Days(){

		$stmt = Connection::connect()->prepare("SELECT DATE(fecha) as sale_date, COUNT(*) as total_count, COALESCE(SUM(total), 0) as total_amount FROM ventas WHERE fecha >= DATE_SUB(CURDATE(), INTERVAL 6 DAY) AND estado = 'completada' GROUP BY DATE(fecha) ORDER BY sale_date ASC");

		$stmt -> execute();

		return $stmt -> fetchAll();

	}

	/*=============================================
	TOP 5 PRODUCTS TODAY (for horizontal bar chart)
	=============================================*/

	static public function mdlTopProductsToday(){

		$stmt = Connection::connect()->prepare("SELECT p.descripcion, SUM(d.cantidad) as total_qty FROM detalle_ventas d INNER JOIN ventas v ON d.id_venta = v.id INNER JOIN productos p ON d.id_producto = p.id WHERE DATE(v.fecha) = CURDATE() AND v.estado = 'completada' GROUP BY d.id_producto ORDER BY total_qty DESC LIMIT 5");

		$stmt -> execute();

		return $stmt -> fetchAll();

	}

	/*=============================================
	SALES BY PAYMENT METHOD (for doughnut chart)
	=============================================*/

	static public function mdlSalesByPaymentMethod(){

		$stmt = Connection::connect()->prepare("SELECT metodo_pago, COUNT(*) as total_count, COALESCE(SUM(total), 0) as total_amount FROM ventas WHERE DATE(fecha) = CURDATE() AND estado = 'completada' GROUP BY metodo_pago");

		$stmt -> execute();

		return $stmt -> fetchAll();

	}

}
