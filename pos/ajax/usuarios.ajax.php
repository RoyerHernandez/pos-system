<?php

session_start();

if(!isset($_SESSION["loggedIn"]) || $_SESSION["loggedIn"] != "ok"){
	echo json_encode(array("error" => "Unauthorized"));
	return;
}

require_once "../models/connection.php";
require_once "../models/users.model.php";

if(isset($_POST["idUsuario"])){

	$item = "id";
	$value = $_POST["idUsuario"];

	$response = UserModel::mdlShowUsers("usuarios", $item, $value);

	echo json_encode($response);

}
