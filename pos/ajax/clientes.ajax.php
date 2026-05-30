<?php

session_start();

if(!isset($_SESSION["loggedIn"]) || $_SESSION["loggedIn"] != "ok"){
	echo json_encode(array("error" => "Unauthorized"));
	return;
}

require_once "../models/connection.php";
require_once "../models/clients.model.php";

if(isset($_POST["idCliente"])){

	$item = "id";
	$value = $_POST["idCliente"];

	$response = ClientModel::mdlShowClients("clientes", $item, $value);

	echo json_encode($response);

}
