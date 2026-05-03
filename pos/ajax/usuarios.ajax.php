<?php

session_start();

if(!isset($_SESSION["iniciarSesion"]) || $_SESSION["iniciarSesion"] != "ok"){
	echo json_encode(array("error" => "No autorizado"));
	return;
}

require_once "../modelos/conexion.php";
require_once "../modelos/usuarios.modelo.php";

if(isset($_POST["idUsuario"])){

	$item = "id";
	$valor = $_POST["idUsuario"];

	$respuesta = ModeloUsuarios::mdlMostrarUsuarios("usuarios", $item, $valor);

	echo json_encode($respuesta);

}
