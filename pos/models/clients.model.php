<?php

require_once "connection.php";

class ClientModel{

	/*=============================================
	SHOW CLIENTS
	=============================================*/

	static public function mdlShowClients($table, $item, $value){

		if($item != null){

			$stmt = Connection::connect()->prepare("SELECT * FROM $table WHERE $item = :$item");

			$stmt -> bindParam(":".$item, $value, PDO::PARAM_STR);

			$stmt -> execute();

			return $stmt -> fetch();

		}else{

			$stmt = Connection::connect()->prepare("SELECT * FROM $table ORDER BY id DESC");

			$stmt -> execute();

			return $stmt -> fetchAll();

		}

	}

	/*=============================================
	SEARCH CLIENTS (autocomplete for POS)
	=============================================*/

	static public function mdlSearchClients($table, $search){

		$search = "%".$search."%";

		$stmt = Connection::connect()->prepare("SELECT id, nombre, documento, telefono, email FROM $table WHERE estado = 1 AND (nombre LIKE :search OR documento LIKE :search2) ORDER BY nombre ASC LIMIT 10");

		$stmt -> bindParam(":search", $search, PDO::PARAM_STR);
		$stmt -> bindParam(":search2", $search, PDO::PARAM_STR);

		$stmt -> execute();

		return $stmt -> fetchAll();

	}

	/*=============================================
	INSERT CLIENT
	=============================================*/

	static public function mdlInsertClient($table, $data){

		$stmt = Connection::connect()->prepare("INSERT INTO $table(nombre, documento, email, telefono, direccion, fecha_nacimiento) VALUES (:nombre, :documento, :email, :telefono, :direccion, :fecha_nacimiento)");

		$stmt -> bindParam(":nombre", $data["nombre"], PDO::PARAM_STR);
		$stmt -> bindParam(":documento", $data["documento"], PDO::PARAM_STR);
		$stmt -> bindParam(":email", $data["email"], PDO::PARAM_STR);
		$stmt -> bindParam(":telefono", $data["telefono"], PDO::PARAM_STR);
		$stmt -> bindParam(":direccion", $data["direccion"], PDO::PARAM_STR);
		$stmt -> bindParam(":fecha_nacimiento", $data["fecha_nacimiento"], PDO::PARAM_STR);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	UPDATE CLIENT
	=============================================*/

	static public function mdlUpdateClient($table, $data){

		$stmt = Connection::connect()->prepare("UPDATE $table SET nombre = :nombre, documento = :documento, email = :email, telefono = :telefono, direccion = :direccion, fecha_nacimiento = :fecha_nacimiento WHERE id = :id");

		$stmt -> bindParam(":id", $data["id"], PDO::PARAM_INT);
		$stmt -> bindParam(":nombre", $data["nombre"], PDO::PARAM_STR);
		$stmt -> bindParam(":documento", $data["documento"], PDO::PARAM_STR);
		$stmt -> bindParam(":email", $data["email"], PDO::PARAM_STR);
		$stmt -> bindParam(":telefono", $data["telefono"], PDO::PARAM_STR);
		$stmt -> bindParam(":direccion", $data["direccion"], PDO::PARAM_STR);
		$stmt -> bindParam(":fecha_nacimiento", $data["fecha_nacimiento"], PDO::PARAM_STR);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	DELETE CLIENT (soft delete)
	=============================================*/

	static public function mdlDeleteClient($table, $data){

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
	UPDATE PURCHASES
	=============================================*/

	static public function mdlUpdatePurchases($table, $id, $amount){

		$stmt = Connection::connect()->prepare("UPDATE $table SET total_compras = total_compras + :amount WHERE id = :id");

		$stmt -> bindParam(":id", $id, PDO::PARAM_INT);
		$stmt -> bindParam(":amount", $amount, PDO::PARAM_STR);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

}
