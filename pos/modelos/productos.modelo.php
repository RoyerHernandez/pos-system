<?php

require_once "conexion.php";

class ModeloProductos{

	/*=============================================
	MOSTRAR PRODUCTOS
	=============================================*/

	static public function mdlMostrarProductos($tabla, $item, $valor){

		if($item != null){

			$stmt = Conexion::conectar()->prepare("SELECT p.*, c.nombre as categoria FROM $tabla p LEFT JOIN categorias c ON p.id_categoria = c.id WHERE p.$item = :$item");

			$stmt -> bindParam(":".$item, $valor, PDO::PARAM_STR);

			$stmt -> execute();

			return $stmt -> fetch();

		}else{

			$stmt = Conexion::conectar()->prepare("SELECT p.*, c.nombre as categoria FROM $tabla p LEFT JOIN categorias c ON p.id_categoria = c.id ORDER BY p.id DESC");

			$stmt -> execute();

			return $stmt -> fetchAll();

		}

	}

	/*=============================================
	INGRESAR PRODUCTO
	=============================================*/

	static public function mdlIngresarProducto($tabla, $datos){

		$stmt = Conexion::conectar()->prepare("INSERT INTO $tabla(codigo, codigo_barras, descripcion, id_categoria, precio_compra, precio_venta, stock, stock_minimo, imagen) VALUES (:codigo, :codigo_barras, :descripcion, :id_categoria, :precio_compra, :precio_venta, :stock, :stock_minimo, :imagen)");

		$stmt -> bindParam(":codigo", $datos["codigo"], PDO::PARAM_STR);
		$stmt -> bindParam(":codigo_barras", $datos["codigo_barras"], PDO::PARAM_STR);
		$stmt -> bindParam(":descripcion", $datos["descripcion"], PDO::PARAM_STR);
		$stmt -> bindParam(":id_categoria", $datos["id_categoria"], PDO::PARAM_INT);
		$stmt -> bindParam(":precio_compra", $datos["precio_compra"], PDO::PARAM_STR);
		$stmt -> bindParam(":precio_venta", $datos["precio_venta"], PDO::PARAM_STR);
		$stmt -> bindParam(":stock", $datos["stock"], PDO::PARAM_INT);
		$stmt -> bindParam(":stock_minimo", $datos["stock_minimo"], PDO::PARAM_INT);
		$stmt -> bindParam(":imagen", $datos["imagen"], PDO::PARAM_STR);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	EDITAR PRODUCTO
	=============================================*/

	static public function mdlEditarProducto($tabla, $datos){

		$stmt = Conexion::conectar()->prepare("UPDATE $tabla SET codigo = :codigo, codigo_barras = :codigo_barras, descripcion = :descripcion, id_categoria = :id_categoria, precio_compra = :precio_compra, precio_venta = :precio_venta, stock = :stock, stock_minimo = :stock_minimo, imagen = :imagen WHERE id = :id");

		$stmt -> bindParam(":id", $datos["id"], PDO::PARAM_INT);
		$stmt -> bindParam(":codigo", $datos["codigo"], PDO::PARAM_STR);
		$stmt -> bindParam(":codigo_barras", $datos["codigo_barras"], PDO::PARAM_STR);
		$stmt -> bindParam(":descripcion", $datos["descripcion"], PDO::PARAM_STR);
		$stmt -> bindParam(":id_categoria", $datos["id_categoria"], PDO::PARAM_INT);
		$stmt -> bindParam(":precio_compra", $datos["precio_compra"], PDO::PARAM_STR);
		$stmt -> bindParam(":precio_venta", $datos["precio_venta"], PDO::PARAM_STR);
		$stmt -> bindParam(":stock", $datos["stock"], PDO::PARAM_INT);
		$stmt -> bindParam(":stock_minimo", $datos["stock_minimo"], PDO::PARAM_INT);
		$stmt -> bindParam(":imagen", $datos["imagen"], PDO::PARAM_STR);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	BORRAR PRODUCTO (soft delete)
	=============================================*/

	static public function mdlBorrarProducto($tabla, $datos){

		$stmt = Conexion::conectar()->prepare("UPDATE $tabla SET estado = :estado WHERE id = :id");

		$stmt -> bindParam(":id", $datos["id"], PDO::PARAM_INT);
		$stmt -> bindParam(":estado", $datos["estado"], PDO::PARAM_INT);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	ACTUALIZAR STOCK
	=============================================*/

	static public function mdlActualizarStock($tabla, $id, $cantidad, $operacion){

		if($operacion == "sumar"){

			$stmt = Conexion::conectar()->prepare("UPDATE $tabla SET stock = stock + :cantidad WHERE id = :id");

		}else{

			$stmt = Conexion::conectar()->prepare("UPDATE $tabla SET stock = stock - :cantidad WHERE id = :id");

		}

		$stmt -> bindParam(":id", $id, PDO::PARAM_INT);
		$stmt -> bindParam(":cantidad", $cantidad, PDO::PARAM_INT);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

}
