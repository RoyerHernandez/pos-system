<?php

require_once "connection.php";

class SaleModel{

	/*=============================================
	INSERT SALE
	=============================================*/

	static public function mdlInsertSale($table, $data){

		$dbh = Connection::connect();

		$stmt = $dbh->prepare("INSERT INTO $table(id_usuario, id_cliente, id_caja, codigo_venta, subtotal, impuesto, descuento, total, metodo_pago) VALUES (:id_usuario, :id_cliente, :id_caja, :codigo_venta, :subtotal, :impuesto, :descuento, :total, :metodo_pago)");

		$stmt -> bindParam(":id_usuario", $data["id_usuario"], PDO::PARAM_INT);
		$stmt -> bindParam(":id_cliente", $data["id_cliente"], PDO::PARAM_INT);
		$stmt -> bindParam(":id_caja", $data["id_caja"], PDO::PARAM_INT);
		$stmt -> bindParam(":codigo_venta", $data["codigo_venta"], PDO::PARAM_STR);
		$stmt -> bindParam(":subtotal", $data["subtotal"], PDO::PARAM_STR);
		$stmt -> bindParam(":impuesto", $data["impuesto"], PDO::PARAM_STR);
		$stmt -> bindParam(":descuento", $data["descuento"], PDO::PARAM_STR);
		$stmt -> bindParam(":total", $data["total"], PDO::PARAM_STR);
		$stmt -> bindParam(":metodo_pago", $data["metodo_pago"], PDO::PARAM_STR);

		if($stmt -> execute()){

			return $dbh->lastInsertId();

		}else{

			return "error";

		}

	}

	/*=============================================
	INSERT SALE DETAIL
	=============================================*/

	static public function mdlInsertSaleDetail($table, $data){

		$stmt = Connection::connect()->prepare("INSERT INTO $table(id_venta, id_producto, cantidad, precio_unitario, descuento, subtotal) VALUES (:id_venta, :id_producto, :cantidad, :precio_unitario, :descuento, :subtotal)");

		$stmt -> bindParam(":id_venta", $data["id_venta"], PDO::PARAM_INT);
		$stmt -> bindParam(":id_producto", $data["id_producto"], PDO::PARAM_INT);
		$stmt -> bindParam(":cantidad", $data["cantidad"], PDO::PARAM_INT);
		$stmt -> bindParam(":precio_unitario", $data["precio_unitario"], PDO::PARAM_STR);
		$stmt -> bindParam(":descuento", $data["descuento"], PDO::PARAM_STR);
		$stmt -> bindParam(":subtotal", $data["subtotal"], PDO::PARAM_STR);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	SHOW SALES
	=============================================*/

	static public function mdlShowSales($table, $item, $value){

		if($item != null){

			$stmt = Connection::connect()->prepare("SELECT v.*, u.nombre as vendedor, c.nombre as cliente FROM $table v LEFT JOIN usuarios u ON v.id_usuario = u.id LEFT JOIN clientes c ON v.id_cliente = c.id WHERE v.$item = :$item");

			$stmt -> bindParam(":".$item, $value, PDO::PARAM_STR);

			$stmt -> execute();

			return $stmt -> fetch();

		}else{

			$stmt = Connection::connect()->prepare("SELECT v.*, u.nombre as vendedor, c.nombre as cliente FROM $table v LEFT JOIN usuarios u ON v.id_usuario = u.id LEFT JOIN clientes c ON v.id_cliente = c.id ORDER BY v.id DESC");

			$stmt -> execute();

			return $stmt -> fetchAll();

		}

	}

	/*=============================================
	SHOW SALE DETAILS
	=============================================*/

	static public function mdlShowSaleDetails($table, $idSale){

		$stmt = Connection::connect()->prepare("SELECT d.*, p.descripcion as producto, p.codigo FROM $table d LEFT JOIN productos p ON d.id_producto = p.id WHERE d.id_venta = :id_venta");

		$stmt -> bindParam(":id_venta", $idSale, PDO::PARAM_INT);

		$stmt -> execute();

		return $stmt -> fetchAll();

	}

	/*=============================================
	CANCEL SALE
	=============================================*/

	static public function mdlCancelSale($table, $id){

		$stmt = Connection::connect()->prepare("UPDATE $table SET estado = 'cancelada' WHERE id = :id");

		$stmt -> bindParam(":id", $id, PDO::PARAM_INT);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	SHOW SALES BY DATE RANGE
	=============================================*/

	static public function mdlShowSalesByDate($table, $startDate, $endDate){

		$stmt = Connection::connect()->prepare("SELECT v.*, u.nombre as vendedor, c.nombre as cliente FROM $table v LEFT JOIN usuarios u ON v.id_usuario = u.id LEFT JOIN clientes c ON v.id_cliente = c.id WHERE DATE(v.fecha) BETWEEN :startDate AND :endDate ORDER BY v.id DESC");

		$stmt -> bindParam(":startDate", $startDate, PDO::PARAM_STR);
		$stmt -> bindParam(":endDate", $endDate, PDO::PARAM_STR);

		$stmt -> execute();

		return $stmt -> fetchAll();

	}

	/*=============================================
	GENERATE SALE CODE
	=============================================*/

	static public function mdlGenerateSaleCode($table){

		$stmt = Connection::connect()->prepare("SELECT codigo_venta FROM $table ORDER BY id DESC LIMIT 1");

		$stmt -> execute();

		$result = $stmt -> fetch();

		if($result){

			$number = intval(substr($result["codigo_venta"], 4)) + 1;

		}else{

			$number = 1;

		}

		return "VTA-" . str_pad($number, 6, "0", STR_PAD_LEFT);

	}

}
