<?php

class ControladorClientes{

	/*=============================================
	MOSTRAR CLIENTES
	=============================================*/

	static public function ctrMostrarClientes($item, $valor){

		$tabla = "clientes";

		$respuesta = ModeloClientes::mdlMostrarClientes($tabla, $item, $valor);

		return $respuesta;

	}

	/*=============================================
	CREAR CLIENTE
	=============================================*/

	static public function ctrCrearCliente(){

		if(isset($_POST["nuevoCliente"])){

			if(preg_match('/^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ ]+$/', $_POST["nuevoCliente"])){

				$tabla = "clientes";

				$fechaNacimiento = !empty($_POST["nuevaFechaNacimiento"]) ? $_POST["nuevaFechaNacimiento"] : null;

				$datos = array("nombre" => $_POST["nuevoCliente"],
					           "documento" => $_POST["nuevoDocumento"],
					           "email" => $_POST["nuevoEmail"],
					           "telefono" => $_POST["nuevoTelefono"],
					           "direccion" => $_POST["nuevaDireccion"],
					           "fecha_nacimiento" => $fechaNacimiento);

				$respuesta = ModeloClientes::mdlIngresarCliente($tabla, $datos);

				if($respuesta == "ok"){

					echo '<script>

					swal({

						type: "success",
						title: "¡El cliente ha sido guardado correctamente!",
						showConfirmButton: true,
						confirmButtonText: "Cerrar"

					}).then(function(result){

						if(result.value){

							window.location = "clientes";

						}

					});

					</script>';

				}

			}else{

				echo '<script>

				swal({

					type: "error",
					title: "¡El nombre no puede ir vacío o llevar caracteres especiales!",
					showConfirmButton: true,
					confirmButtonText: "Cerrar"

				}).then(function(result){

					if(result.value){

						window.location = "clientes";

					}

				});

				</script>';

			}

		}

	}

	/*=============================================
	EDITAR CLIENTE
	=============================================*/

	static public function ctrEditarCliente(){

		if(isset($_POST["editarCliente"])){

			if(preg_match('/^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ ]+$/', $_POST["editarCliente"])){

				$tabla = "clientes";

				$fechaNacimiento = !empty($_POST["editarFechaNacimiento"]) ? $_POST["editarFechaNacimiento"] : null;

				$datos = array("id" => $_POST["idClienteEditar"],
					           "nombre" => $_POST["editarCliente"],
					           "documento" => $_POST["editarDocumento"],
					           "email" => $_POST["editarEmail"],
					           "telefono" => $_POST["editarTelefono"],
					           "direccion" => $_POST["editarDireccion"],
					           "fecha_nacimiento" => $fechaNacimiento);

				$respuesta = ModeloClientes::mdlEditarCliente($tabla, $datos);

				if($respuesta == "ok"){

					echo '<script>

					swal({

						type: "success",
						title: "¡El cliente ha sido editado correctamente!",
						showConfirmButton: true,
						confirmButtonText: "Cerrar"

					}).then(function(result){

						if(result.value){

							window.location = "clientes";

						}

					});

					</script>';

				}

			}else{

				echo '<script>

				swal({

					type: "error",
					title: "¡El nombre no puede ir vacío o llevar caracteres especiales!",
					showConfirmButton: true,
					confirmButtonText: "Cerrar"

				}).then(function(result){

					if(result.value){

						window.location = "clientes";

					}

				});

				</script>';

			}

		}

	}

	/*=============================================
	BORRAR CLIENTE
	=============================================*/

	static public function ctrBorrarCliente(){

		if(isset($_GET["idCliente"])){

			$tabla = "clientes";

			$datos = array("id" => $_GET["idCliente"],
				           "estado" => 0);

			$respuesta = ModeloClientes::mdlBorrarCliente($tabla, $datos);

			if($respuesta == "ok"){

				echo '<script>

				swal({

					type: "success",
					title: "¡El cliente ha sido eliminado correctamente!",
					showConfirmButton: true,
					confirmButtonText: "Cerrar"

				}).then(function(result){

					if(result.value){

						window.location = "clientes";

					}

				});

				</script>';

			}

		}

	}

}
