<?php

session_start();

if(!isset($_SESSION["loggedIn"]) || $_SESSION["loggedIn"] != "ok"){
	echo json_encode(array("status" => "error", "message" => "No autorizado"));
	return;
}

require_once "../models/connection.php";
require_once "../models/inventory.model.php";
require_once "../models/products.model.php";

require_once "../controllers/inventory.controller.php";

/*=============================================
CREATE MOVEMENT
=============================================*/

if(isset($_POST["createMovement"])){

	$data = array(
		"id_producto" => $_POST["idProducto"],
		"id_usuario" => $_SESSION["id"],
		"tipo" => $_POST["tipo"],
		"motivo" => $_POST["motivo"],
		"cantidad" => $_POST["cantidad"],
		"observaciones" => isset($_POST["observaciones"]) ? $_POST["observaciones"] : null
	);

	$response = InventoryController::ctrCreateMovement($data);

	echo json_encode($response);

}

/*=============================================
GET MOVEMENTS
=============================================*/

if(isset($_POST["getMovements"])){

	if(isset($_POST["startDate"]) && isset($_POST["endDate"]) && $_POST["startDate"] != "" && $_POST["endDate"] != ""){

		$response = InventoryController::ctrShowMovementsByDate($_POST["startDate"], $_POST["endDate"]);

	}else{

		$response = InventoryController::ctrShowMovements(null, null);

	}

	echo json_encode($response);

}

/*=============================================
GET MOTIVOS
=============================================*/

if(isset($_POST["getMotivos"])){

	$response = InventoryController::ctrGetMotivos();

	echo json_encode($response);

}
