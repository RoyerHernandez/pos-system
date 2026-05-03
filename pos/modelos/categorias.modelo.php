<?php

require_once "conexion.php";

class CategoryModel{

	/*=============================================
	SHOW CATEGORIES
	=============================================*/

	static public function mdlShowCategories($table, $item, $value){

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
	INSERT CATEGORY
	=============================================*/

	static public function mdlInsertCategory($table, $data){

		$stmt = Connection::connect()->prepare("INSERT INTO $table(nombre, descripcion) VALUES (:nombre, :descripcion)");

		$stmt -> bindParam(":nombre", $data["nombre"], PDO::PARAM_STR);
		$stmt -> bindParam(":descripcion", $data["descripcion"], PDO::PARAM_STR);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	UPDATE CATEGORY
	=============================================*/

	static public function mdlUpdateCategory($table, $data){

		$stmt = Connection::connect()->prepare("UPDATE $table SET nombre = :nombre, descripcion = :descripcion WHERE id = :id");

		$stmt -> bindParam(":id", $data["id"], PDO::PARAM_INT);
		$stmt -> bindParam(":nombre", $data["nombre"], PDO::PARAM_STR);
		$stmt -> bindParam(":descripcion", $data["descripcion"], PDO::PARAM_STR);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

	/*=============================================
	DELETE CATEGORY (soft delete)
	=============================================*/

	static public function mdlDeleteCategory($table, $data){

		$stmt = Connection::connect()->prepare("UPDATE $table SET estado = :estado WHERE id = :id");

		$stmt -> bindParam(":id", $data["id"], PDO::PARAM_INT);
		$stmt -> bindParam(":estado", $data["estado"], PDO::PARAM_INT);

		if($stmt -> execute()){

			return "ok";

		}else{

			return "error";

		}

	}

}
