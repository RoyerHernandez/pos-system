<?php

require_once "connection.php";

class ProductModel{

	/*=============================================
	SHOW PRODUCTS
	=============================================*/

	static public function mdlShowProducts($table, $item, $value){

		if($item != null){

			$stmt = Connection::connect()->prepare("SELECT p.*, c.nombre as categoria FROM $table p LEFT JOIN categorias c ON p.id_categoria = c.id WHERE p.$item = :$item");

			$stmt -> bindParam(":".$item, $value, PDO::PARAM_STR);

			$stmt -> execute();

			return $stmt -> fetch();

		}else{

			$stmt = Connection::connect()->prepare("SELECT p.*, c.nombre as categoria FROM $table p LEFT JOIN categorias c ON p.id_categoria = c.id ORDER BY p.id DESC");

			$stmt -> execute();

			return $stmt -> fetchAll();

		}

	}

	/*=============================================
	INSERT PRODUCT
	=============================================*/

	static public function mdlInsertProduct($table, $data){

		$stmt = Connection::connect()->prepare("INSERT INTO $table(codigo, codigo_barras, descripcion, id_categoria, precio_compra, precio_venta, stock, stock_minimo, imagen) VALUES (:codigo, :codigo_barras, :descripcion, :id_categoria, :precio_compra, :precio_venta, :stock, :stock_minimo, :imagen)");

		$stmt -> bindParam(":codigo", $data["codigo"], PDO::PARAM_STR);
		$stmt -> bindParam(":codigo_barras", $data["codigo_barras"], PDO::PARAM_STR);
		$stmt -> bindParam(":descripcion", $data["descripcion"], PDO::PARAM_STR);
		$stmt -> bindParam(":id_categoria", $data["id_categoria"], PDO::PARAM_INT);
		$stmt -> bindParam(":precio_compra", $data["precio_compra"], PDO::PARAM_STR);
		$stmt -> bindParam(":precio_venta", $data["precio_venta"], PDO::PARAM_STR);
		$stmt -> bindParam(":stock", $data["stock"], PDO::PARAM_INT);
		$stmt -> bindParam(":stock_minimo", $data["stock_minimo"], PDO::PARAM_INT);
		$stmt -> bindParam(":imagen", $data["imagen"], PDO::PARAM_STR);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	UPDATE PRODUCT
	=============================================*/

	static public function mdlUpdateProduct($table, $data){

		$stmt = Connection::connect()->prepare("UPDATE $table SET codigo = :codigo, codigo_barras = :codigo_barras, descripcion = :descripcion, id_categoria = :id_categoria, precio_compra = :precio_compra, precio_venta = :precio_venta, stock = :stock, stock_minimo = :stock_minimo, imagen = :imagen WHERE id = :id");

		$stmt -> bindParam(":id", $data["id"], PDO::PARAM_INT);
		$stmt -> bindParam(":codigo", $data["codigo"], PDO::PARAM_STR);
		$stmt -> bindParam(":codigo_barras", $data["codigo_barras"], PDO::PARAM_STR);
		$stmt -> bindParam(":descripcion", $data["descripcion"], PDO::PARAM_STR);
		$stmt -> bindParam(":id_categoria", $data["id_categoria"], PDO::PARAM_INT);
		$stmt -> bindParam(":precio_compra", $data["precio_compra"], PDO::PARAM_STR);
		$stmt -> bindParam(":precio_venta", $data["precio_venta"], PDO::PARAM_STR);
		$stmt -> bindParam(":stock", $data["stock"], PDO::PARAM_INT);
		$stmt -> bindParam(":stock_minimo", $data["stock_minimo"], PDO::PARAM_INT);
		$stmt -> bindParam(":imagen", $data["imagen"], PDO::PARAM_STR);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	DELETE PRODUCT (soft delete)
	=============================================*/

	static public function mdlDeleteProduct($table, $data){

		$stmt = Connection::connect()->prepare("UPDATE $table SET estado = :estado WHERE id = :id");

		$stmt -> bindParam(":id", $data["id"], PDO::PARAM_INT);
		$stmt -> bindParam(":estado", $data["estado"], PDO::PARAM_INT);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	SEARCH PRODUCTS (autocomplete for POS)
	=============================================*/

	static public function mdlSearchProducts($table, $search){

		$search = "%".$search."%";

		$stmt = Connection::connect()->prepare("SELECT p.id, p.codigo, p.codigo_barras, p.descripcion, p.precio_venta, p.stock, p.imagen, c.nombre as categoria FROM $table p LEFT JOIN categorias c ON p.id_categoria = c.id WHERE p.estado = 1 AND (p.descripcion LIKE :search OR p.codigo LIKE :search2 OR p.codigo_barras LIKE :search3) LIMIT 10");

		$stmt -> bindParam(":search", $search, PDO::PARAM_STR);
		$stmt -> bindParam(":search2", $search, PDO::PARAM_STR);
		$stmt -> bindParam(":search3", $search, PDO::PARAM_STR);

		$stmt -> execute();

		return $stmt -> fetchAll();

	}

	/*=============================================
	GET PRODUCT BY CODE (barcode scanner)
	=============================================*/

	static public function mdlGetProductByCode($table, $code){

		$stmt = Connection::connect()->prepare("SELECT p.*, c.nombre as categoria FROM $table p LEFT JOIN categorias c ON p.id_categoria = c.id WHERE p.estado = 1 AND (p.codigo = :code OR p.codigo_barras = :code2) LIMIT 1");

		$stmt -> bindParam(":code", $code, PDO::PARAM_STR);
		$stmt -> bindParam(":code2", $code, PDO::PARAM_STR);

		$stmt -> execute();

		return $stmt -> fetch();

	}

	/*=============================================
	UPDATE STOCK
	=============================================*/

	static public function mdlUpdateStock($table, $id, $quantity, $operation){

		if($operation == "add"){

			$stmt = Connection::connect()->prepare("UPDATE $table SET stock = stock + :quantity WHERE id = :id");

		}else{

			$stmt = Connection::connect()->prepare("UPDATE $table SET stock = stock - :quantity WHERE id = :id");

		}

		$stmt -> bindParam(":id", $id, PDO::PARAM_INT);
		$stmt -> bindParam(":quantity", $quantity, PDO::PARAM_INT);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

}
