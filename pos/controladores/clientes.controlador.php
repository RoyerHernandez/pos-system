<?php

class ClientController{

	/*=============================================
	SHOW CLIENTS
	=============================================*/

	static public function ctrShowClients($item, $valor){

		$table = "clientes";

		$response = ClientModel::mdlShowClients($table, $item, $valor);

		return $response;

	}

	/*=============================================
	CREATE CLIENT
	=============================================*/

	static public function ctrCreateClient(){

		if(isset($_POST["nuevoCliente"])){

			if(preg_match('/^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ ]+$/', $_POST["nuevoCliente"])){

				$table = "clientes";

				$birthDate = !empty($_POST["nuevaFechaNacimiento"]) ? $_POST["nuevaFechaNacimiento"] : null;

				$data = array("nombre" => $_POST["nuevoCliente"],
					           "documento" => $_POST["nuevoDocumento"],
					           "email" => $_POST["nuevoEmail"],
					           "telefono" => $_POST["nuevoTelefono"],
					           "direccion" => $_POST["nuevaDireccion"],
					           "fecha_nacimiento" => $birthDate);

				$response = ClientModel::mdlInsertClient($table, $data);

				if($response == "ok"){

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
	UPDATE CLIENT
	=============================================*/

	static public function ctrUpdateClient(){

		if(isset($_POST["editarCliente"])){

			if(preg_match('/^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ ]+$/', $_POST["editarCliente"])){

				$table = "clientes";

				$birthDate = !empty($_POST["editarFechaNacimiento"]) ? $_POST["editarFechaNacimiento"] : null;

				$data = array("id" => $_POST["idClienteEditar"],
					           "nombre" => $_POST["editarCliente"],
					           "documento" => $_POST["editarDocumento"],
					           "email" => $_POST["editarEmail"],
					           "telefono" => $_POST["editarTelefono"],
					           "direccion" => $_POST["editarDireccion"],
					           "fecha_nacimiento" => $birthDate);

				$response = ClientModel::mdlUpdateClient($table, $data);

				if($response == "ok"){

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
	DELETE CLIENT
	=============================================*/

	static public function ctrDeleteClient(){

		if(isset($_GET["idCliente"])){

			$table = "clientes";

			$data = array("id" => $_GET["idCliente"],
				           "estado" => 0);

			$response = ClientModel::mdlDeleteClient($table, $data);

			if($response == "ok"){

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