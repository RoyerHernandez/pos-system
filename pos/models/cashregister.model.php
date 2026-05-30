<?php

require_once "connection.php";

class CashRegisterModel{

	/*=============================================
	OPEN CASH REGISTER
	=============================================*/

	static public function mdlOpenCashRegister($table, $data){

		$stmt = Connection::connect()->prepare("INSERT INTO $table(id_usuario, monto_apertura) VALUES (:id_usuario, :monto_apertura)");

		$stmt -> bindParam(":id_usuario", $data["id_usuario"], PDO::PARAM_INT);
		$stmt -> bindParam(":monto_apertura", $data["monto_apertura"], PDO::PARAM_STR);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	CLOSE CASH REGISTER
	=============================================*/

	static public function mdlCloseCashRegister($table, $data){

		$stmt = Connection::connect()->prepare("UPDATE $table SET monto_cierre = :monto_cierre, total_ventas = :total_ventas, total_efectivo = :total_efectivo, total_tarjeta = :total_tarjeta, total_transferencia = :total_transferencia, estado = 'cerrada', fecha_cierre = NOW() WHERE id = :id");

		$stmt -> bindParam(":id", $data["id"], PDO::PARAM_INT);
		$stmt -> bindParam(":monto_cierre", $data["monto_cierre"], PDO::PARAM_STR);
		$stmt -> bindParam(":total_ventas", $data["total_ventas"], PDO::PARAM_STR);
		$stmt -> bindParam(":total_efectivo", $data["total_efectivo"], PDO::PARAM_STR);
		$stmt -> bindParam(":total_tarjeta", $data["total_tarjeta"], PDO::PARAM_STR);
		$stmt -> bindParam(":total_transferencia", $data["total_transferencia"], PDO::PARAM_STR);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	GET OPEN CASH REGISTER
	=============================================*/

	static public function mdlGetOpenCashRegister($table, $idUsuario){

		$stmt = Connection::connect()->prepare("SELECT * FROM $table WHERE id_usuario = :id_usuario AND estado = 'abierta' LIMIT 1");

		$stmt -> bindParam(":id_usuario", $idUsuario, PDO::PARAM_INT);

		$stmt -> execute();

		return $stmt -> fetch();

	}

	/*=============================================
	UPDATE CASH REGISTER TOTALS
	=============================================*/

	static public function mdlUpdateCashRegisterTotals($table, $id, $total, $paymentMethod){

		$column = "total_efectivo";

		if($paymentMethod == "Tarjeta"){
			$column = "total_tarjeta";
		}else if($paymentMethod == "Transferencia"){
			$column = "total_transferencia";
		}

		$stmt = Connection::connect()->prepare("UPDATE $table SET total_ventas = total_ventas + :total, $column = $column + :total2 WHERE id = :id");

		$stmt -> bindParam(":id", $id, PDO::PARAM_INT);
		$stmt -> bindParam(":total", $total, PDO::PARAM_STR);
		$stmt -> bindParam(":total2", $total, PDO::PARAM_STR);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	GET CASH REGISTER BY ID
	=============================================*/

	static public function mdlGetCashRegister($table, $id){

		$stmt = Connection::connect()->prepare("SELECT c.*, u.nombre as usuario FROM $table c LEFT JOIN usuarios u ON c.id_usuario = u.id WHERE c.id = :id");

		$stmt -> bindParam(":id", $id, PDO::PARAM_INT);

		$stmt -> execute();

		return $stmt -> fetch();

	}

}
