<?php

require_once "connection.php";

class InventoryModel{

	/*=============================================
	INSERT MOVEMENT
	=============================================*/

	static public function mdlInsertMovement($table, $data){

		$stmt = Connection::connect()->prepare("INSERT INTO $table(id_producto, id_usuario, tipo, motivo, cantidad, observaciones, id_referencia) VALUES (:id_producto, :id_usuario, :tipo, :motivo, :cantidad, :observaciones, :id_referencia)");

		$stmt -> bindParam(":id_producto", $data["id_producto"], PDO::PARAM_INT);
		$stmt -> bindParam(":id_usuario", $data["id_usuario"], PDO::PARAM_INT);
		$stmt -> bindParam(":tipo", $data["tipo"], PDO::PARAM_STR);
		$stmt -> bindParam(":motivo", $data["motivo"], PDO::PARAM_STR);
		$stmt -> bindParam(":cantidad", $data["cantidad"], PDO::PARAM_INT);
		$stmt -> bindParam(":observaciones", $data["observaciones"], PDO::PARAM_STR);
		$stmt -> bindParam(":id_referencia", $data["id_referencia"], PDO::PARAM_INT);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	SHOW MOVEMENTS
	=============================================*/

	static public function mdlShowMovements($table, $item, $value){

		if($item != null){

			$stmt = Connection::connect()->prepare("SELECT m.*, p.descripcion as producto, p.codigo as codigo_producto, u.nombre as usuario FROM $table m LEFT JOIN productos p ON m.id_producto = p.id LEFT JOIN usuarios u ON m.id_usuario = u.id WHERE m.$item = :$item ORDER BY m.fecha DESC");

			$stmt -> bindParam(":".$item, $value, PDO::PARAM_STR);

			$stmt -> execute();

			return $stmt -> fetch();

		}else{

			$stmt = Connection::connect()->prepare("SELECT m.*, p.descripcion as producto, p.codigo as codigo_producto, u.nombre as usuario FROM $table m LEFT JOIN productos p ON m.id_producto = p.id LEFT JOIN usuarios u ON m.id_usuario = u.id ORDER BY m.fecha DESC");

			$stmt -> execute();

			return $stmt -> fetchAll();

		}

	}

	/*=============================================
	SHOW MOVEMENTS BY DATE RANGE
	=============================================*/

	static public function mdlShowMovementsByDate($table, $startDate, $endDate){

		$stmt = Connection::connect()->prepare("SELECT m.*, p.descripcion as producto, p.codigo as codigo_producto, u.nombre as usuario FROM $table m LEFT JOIN productos p ON m.id_producto = p.id LEFT JOIN usuarios u ON m.id_usuario = u.id WHERE DATE(m.fecha) BETWEEN :startDate AND :endDate ORDER BY m.fecha DESC");

		$stmt -> bindParam(":startDate", $startDate, PDO::PARAM_STR);
		$stmt -> bindParam(":endDate", $endDate, PDO::PARAM_STR);

		$stmt -> execute();

		return $stmt -> fetchAll();

	}

	/*=============================================
	SHOW MOVEMENTS BY REFERENCE
	=============================================*/

	static public function mdlShowMovementsByReference($table, $idReferencia){

		$stmt = Connection::connect()->prepare("SELECT m.*, p.descripcion as producto, p.codigo as codigo_producto, u.nombre as usuario FROM $table m LEFT JOIN productos p ON m.id_producto = p.id LEFT JOIN usuarios u ON m.id_usuario = u.id WHERE m.id_referencia = :id_referencia ORDER BY m.fecha DESC");

		$stmt -> bindParam(":id_referencia", $idReferencia, PDO::PARAM_INT);

		$stmt -> execute();

		return $stmt -> fetchAll();

	}

}
