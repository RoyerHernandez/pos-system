<?php

class UserController{

	/*=============================================
	USER LOGIN
	=============================================*/

	static public function ctrUserLogin(){

		if(isset($_POST["ingUsuario"])){

			if(preg_match('/^[a-zA-Z0-9]+$/', $_POST["ingUsuario"]) &&
			   preg_match('/^[a-zA-Z0-9]+$/', $_POST["ingPassword"])){

				$table = "usuarios";

				$item = "usuario";
				$value = $_POST["ingUsuario"];

				$response = UserModel::mdlShowUsers($table, $item, $value);

				if($response && $response["usuario"] == $_POST["ingUsuario"] &&
				   password_verify($_POST["ingPassword"], $response["password"]) &&
				   $response["estado"] == 1){

					$_SESSION["loggedIn"] = "ok";
					$_SESSION["id"] = $response["id"];
					$_SESSION["name"] = $response["nombre"];
					$_SESSION["username"] = $response["usuario"];
					$_SESSION["photo"] = $response["foto"];
					$_SESSION["role"] = $response["perfil"];

					UserModel::mdlUpdateLastLogin($table, $response["id"]);

					echo '<script>

						window.location = "inicio";

					</script>';

				}else{

					echo '<br><div class="alert alert-danger">Error al ingresar, vuelve a intentarlo</div>';

				}

			}

		}

	}

	/*=============================================
	SHOW USERS
	=============================================*/

	static public function ctrShowUsers($item, $value){

		$table = "usuarios";

		$response = UserModel::mdlShowUsers($table, $item, $value);

		return $response;

	}

	/*=============================================
	CREATE USER
	=============================================*/

	static public function ctrCreateUser(){

		if(isset($_POST["nuevoUsuario"]) && $_POST["nuevoUsuario"] != ""){

			if(preg_match('/^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ ]+$/', $_POST["nuevoNombre"]) &&
			   preg_match('/^[a-zA-Z0-9]+$/', $_POST["nuevoUsuario"]) &&
			   preg_match('/^[a-zA-Z0-9]+$/', $_POST["nuevoPassword"])){

			   	/*=============================================
				VALIDATE IMAGE
				=============================================*/

				$path = "";

				if(isset($_FILES["nuevaFoto"]["tmp_name"]) && !empty($_FILES["nuevaFoto"]["tmp_name"])){

					list($width, $height) = getimagesize($_FILES["nuevaFoto"]["tmp_name"]);

					$newWidth = 500;
					$newHeight = 500;

					/*=============================================
					CREATE DIRECTORY TO STORE USER PHOTO
					=============================================*/

					$absDir = __DIR__ . "/../views/img/usuarios/" . $_POST["nuevoUsuario"];

					if(!is_dir($absDir)){
						mkdir($absDir, 0755, true);
					}

					/*=============================================
					PROCESS IMAGE (any format supported by GD)
					=============================================*/

					$source = @imagecreatefromstring(file_get_contents($_FILES["nuevaFoto"]["tmp_name"]));

					if($source !== false){

						$random = mt_rand(100,999);

						$path = "views/img/usuarios/".$_POST["nuevoUsuario"]."/".$random.".jpg";

						$destination = imagecreatetruecolor($newWidth, $newHeight);

						imagecopyresized($destination, $source, 0, 0, 0, 0, $newWidth, $newHeight, $width, $height);

						imagejpeg($destination, $absDir."/".$random.".jpg");

						imagedestroy($source);
						imagedestroy($destination);

					}

				}

				$table = "usuarios";

				$encrypted = password_hash($_POST["nuevoPassword"], PASSWORD_BCRYPT);

				$data = array("nombre" => $_POST["nuevoNombre"],
					           "usuario" => $_POST["nuevoUsuario"],
					           "password" => $encrypted,
					           "perfil" => $_POST["nuevoPerfil"],
					           "foto" => $path);

				$response = UserModel::mdlInsertUser($table, $data);

				if($response == "ok"){

					echo '<script>

					swal({

						type: "success",
						title: "¡El usuario ha sido guardado correctamente!",
						showConfirmButton: true,
						confirmButtonText: "Cerrar"

					}).then(function(result){

						if(result.value){

							window.location = "usuarios";

						}

					});


					</script>';


				}


			}else{

				echo '<script>

					swal({

						type: "error",
						title: "¡El usuario no puede ir vacío o llevar caracteres especiales!",
						showConfirmButton: true,
						confirmButtonText: "Cerrar"

					}).then(function(result){

						if(result.value){

							window.location = "usuarios";

						}

					});


				</script>';

			}
		}
	}

	/*=============================================
	UPDATE USER
	=============================================*/

	static public function ctrUpdateUser(){

		if(isset($_POST["editarUsuario"])){

			/*=============================================
			VALIDATE IMAGE
			=============================================*/

			$path = $_POST["fotoActual"];

			if(isset($_FILES["editarFoto"]["tmp_name"]) && !empty($_FILES["editarFoto"]["tmp_name"])){

				list($width, $height) = getimagesize($_FILES["editarFoto"]["tmp_name"]);

				$newWidth = 500;
				$newHeight = 500;

				$absDir = __DIR__ . "/../views/img/usuarios/" . $_POST["editarUsuario"];

				if(!is_dir($absDir)){
					mkdir($absDir, 0755, true);
				}

				/*=============================================
				PROCESS IMAGE (any format supported by GD)
				=============================================*/

				$source = @imagecreatefromstring(file_get_contents($_FILES["editarFoto"]["tmp_name"]));

				if($source !== false){

					$random = mt_rand(100,999);

					$path = "views/img/usuarios/".$_POST["editarUsuario"]."/".$random.".jpg";

					$destination = imagecreatetruecolor($newWidth, $newHeight);

					imagecopyresized($destination, $source, 0, 0, 0, 0, $newWidth, $newHeight, $width, $height);

					imagejpeg($destination, $absDir."/".$random.".jpg");

					imagedestroy($source);
					imagedestroy($destination);

				}

			}

			$table = "usuarios";

			/*=============================================
			CHECK IF PASSWORD WAS CHANGED
			=============================================*/

			if(!empty($_POST["editarPassword"])){
				$password = password_hash($_POST["editarPassword"], PASSWORD_BCRYPT);
			}else{
				$password = $_POST["passwordActual"];
			}

			$data = array("nombre" => $_POST["editarNombre"],
				           "usuario" => $_POST["editarUsuario"],
				           "password" => $password,
				           "perfil" => $_POST["editarPerfil"],
				           "foto" => $path);

			$response = UserModel::mdlUpdateUser($table, $data);

			if($response == "ok"){

				echo '<script>

				swal({

					type: "success",
					title: "¡El usuario ha sido editado correctamente!",
					showConfirmButton: true,
					confirmButtonText: "Cerrar"

				}).then(function(result){

					if(result.value){

						window.location = "usuarios";

					}

				});

				</script>';

			}

		}

	}

	/*=============================================
	DELETE USER
	=============================================*/

	static public function ctrDeleteUser(){

		if(isset($_GET["idUsuario"])){

			$table = "usuarios";

			$data = array("id" => $_GET["idUsuario"],
				           "estado" => 0);

			$response = UserModel::mdlDeleteUser($table, $data);

			if($response == "ok"){

				echo '<script>

				swal({

					type: "success",
					title: "¡El usuario ha sido eliminado correctamente!",
					showConfirmButton: true,
					confirmButtonText: "Cerrar"

				}).then(function(result){

					if(result.value){

						window.location = "usuarios";

					}

				});

				</script>';

			}

		}

	}

}
