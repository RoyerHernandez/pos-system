<?php

session_start();

if(!isset($_SESSION["iniciarSesion"]) || $_SESSION["iniciarSesion"] != "ok"){
	echo json_encode(array("error" => "No autorizado"));
	return;
}

require_once "../modelos/conexion.php";
require_once "../modelos/clientes.modelo.php";

if(isset($_POST["idCliente"])){

	$item = "id";
	$valor = $_POST["idCliente"];

	$respuesta = ModeloClientes::mdlMostrarClientes("clientes", $item, $valor);

	echo json_encode($respuesta);

}
