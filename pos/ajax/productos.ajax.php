<?php

session_start();

if(!isset($_SESSION["loggedIn"]) || $_SESSION["loggedIn"] != "ok"){
	echo json_encode(array("error" => "Unauthorized"));
	return;
}

require_once "../models/connection.php";
require_once "../models/products.model.php";

if(isset($_POST["idProducto"])){

	$item = "p.id";
	$value = $_POST["idProducto"];

	$response = ProductModel::mdlShowProducts("productos", $item, $value);

	echo json_encode($response);

}
