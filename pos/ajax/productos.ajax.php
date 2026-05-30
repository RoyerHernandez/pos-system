<?php

session_start();

if(!isset($_SESSION["loggedIn"]) || $_SESSION["loggedIn"] != "ok"){
	echo json_encode(array("error" => "Unauthorized"));
	return;
}

require_once "../models/connection.php";
require_once "../models/products.model.php";

/*=============================================
GET PRODUCT BY ID (for edit modal)
=============================================*/

if(isset($_POST["idProducto"])){

	$item = "p.id";
	$value = $_POST["idProducto"];

	$response = ProductModel::mdlShowProducts("productos", $item, $value);

	echo json_encode($response);

}

/*=============================================
SEARCH PRODUCTS (autocomplete for POS)
=============================================*/

if(isset($_POST["searchProduct"])){

	$search = $_POST["searchProduct"];

	$response = ProductModel::mdlSearchProducts("productos", $search);

	echo json_encode($response);

}

/*=============================================
GET PRODUCT BY CODE (barcode scanner for POS)
=============================================*/

if(isset($_POST["productCode"])){

	$code = $_POST["productCode"];

	$response = ProductModel::mdlGetProductByCode("productos", $code);

	echo json_encode($response);

}
