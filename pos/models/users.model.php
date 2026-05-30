<?php

require_once "connection.php";

class UserModel{

	/*=============================================
	SHOW USERS
	=============================================*/

	static public function mdlShowUsers($table, $item, $value){

		if($item != null){

			$stmt = Connection::connect()->prepare("SELECT * FROM $table WHERE $item = :$item");

			$stmt -> bindParam(":".$item, $value, PDO::PARAM_STR);

			$stmt -> execute();

			return $stmt -> fetch();

		}else{

			$stmt = Connection::connect()->prepare("SELECT * FROM $table");

			$stmt -> execute();

			return $stmt -> fetchAll();

		}

	}

	/*=============================================
	INSERT USER
	=============================================*/

	static public function mdlInsertUser($table, $data){

		$stmt = Connection::connect()->prepare("INSERT INTO $table(nombre, usuario, password, perfil, foto) VALUES (:nombre, :usuario, :password, :perfil, :foto)");

		$stmt -> bindParam(":nombre", $data["nombre"], PDO::PARAM_STR);
		$stmt -> bindParam(":usuario", $data["usuario"], PDO::PARAM_STR);
		$stmt -> bindParam(":password", $data["password"], PDO::PARAM_STR);
		$stmt -> bindParam(":perfil", $data["perfil"], PDO::PARAM_STR);
		$stmt -> bindParam(":foto", $data["foto"], PDO::PARAM_STR);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	UPDATE USER
	=============================================*/

	static public function mdlUpdateUser($table, $data){

		$stmt = Connection::connect()->prepare("UPDATE $table SET nombre = :nombre, password = :password, perfil = :perfil, foto = :foto WHERE usuario = :usuario");

		$stmt -> bindParam(":nombre", $data["nombre"], PDO::PARAM_STR);
		$stmt -> bindParam(":usuario", $data["usuario"], PDO::PARAM_STR);
		$stmt -> bindParam(":password", $data["password"], PDO::PARAM_STR);
		$stmt -> bindParam(":perfil", $data["perfil"], PDO::PARAM_STR);
		$stmt -> bindParam(":foto", $data["foto"], PDO::PARAM_STR);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	DELETE USER (soft delete)
	=============================================*/

	static public function mdlDeleteUser($table, $data){

		$stmt = Connection::connect()->prepare("UPDATE $table SET estado = :estado WHERE id = :id");

		$stmt -> bindParam(":estado", $data["estado"], PDO::PARAM_INT);
		$stmt -> bindParam(":id", $data["id"], PDO::PARAM_INT);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	UPDATE LAST LOGIN
	=============================================*/

	static public function mdlUpdateLastLogin($table, $id){

		$stmt = Connection::connect()->prepare("UPDATE $table SET ultimo_login = NOW() WHERE id = :id");

		$stmt -> bindParam(":id", $id, PDO::PARAM_INT);

		$stmt -> execute();

	}

}