<?php

session_start();

if(!isset($_SESSION["iniciarSesion"]) || $_SESSION["iniciarSesion"] != "ok"){
	echo json_encode(array("error" => "No autorizado"));
	return;
}

require_once "../modelos/conexion.php";
require_once "../modelos/productos.modelo.php";

if(isset($_POST["idProducto"])){

	$item = "p.id";
	$valor = $_POST["idProducto"];

	$respuesta = ModeloProductos::mdlMostrarProductos("productos", $item, $valor);

	echo json_encode($respuesta);

}
