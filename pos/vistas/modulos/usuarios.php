<?php

$usuarios = ControladorUsuarios::ctrMostrarUsuarios(null, null);

?>

<div class="content-wrapper">

  <section class="content-header">

    <h1>
      Administrar usuarios
    </h1>

    <ol class="breadcrumb">
      <li><a href="inicio"><i class="fa fa-dashboard"></i> Inicio</a></li>
      <li class="active">Administrar usuarios</li>
    </ol>

  </section>

  <section class="content">

    <div class="box">

      <div class="box-header with-border">

        <button class="btn btn-primary" data-toggle="modal" data-target="#modalAgregarUsuario">
          <i class="fa fa-plus"></i> Agregar Usuario
        </button>

      </div>

      <div class="box-body">

        <table class="table table-bordered table-striped dt-responsive tablas">

          <thead>
            <tr>
              <th style="width:10px;">#</th>
              <th>Nombre</th>
              <th>Usuario</th>
              <th>Foto</th>
              <th>Perfil</th>
              <th>Estado</th>
              <th>Último login</th>
              <th>Acciones</th>
            </tr>
          </thead>

          <tbody>

            <?php

            $contador = 1;

            foreach($usuarios as $key => $valor){

              echo '<tr>

                <td>'.$contador.'</td>
                <td>'.$valor["nombre"].'</td>
                <td>'.$valor["usuario"].'</td>
                <td><img src="'.(!empty($valor["foto"]) ? $valor["foto"] : 'vistas/img/usuarios/default/anonymous.png').'" class="img-thumbnail" width="40px"></td>
                <td>'.$valor["perfil"].'</td>
                <td>';

              if($valor["estado"] == 1){
                echo '<button class="btn btn-success btn-xs">Activado</button>';
              }else{
                echo '<button class="btn btn-danger btn-xs">Inactivo</button>';
              }

              echo '</td>
                <td>'.$valor["ultimo_login"].'</td>
                <td>

                  <div class="btn-group">

                    <button class="btn btn-warning btnEditarUsuario" data-toggle="modal" data-target="#modalEditarUsuario" idUsuario="'.$valor["id"].'" ><i class="fa fa-pencil"></i></button>

                    <a class="btn btn-danger btnEliminarUsuario" idUsuario="'.$valor["id"].'" usuario="'.$valor["usuario"].'"><i class="fa fa-times"></i></a>

                  </div>

                </td>

              </tr>';

              $contador++;

            }

            ?>

          </tbody>

        </table>

      </div>

    </div>

  </section>

</div>

<!--=====================================
MODAL AGREGAR USUARIO
======================================-->

<div class="modal fade" id="modalAgregarUsuario">
  <div class="modal-dialog">
    <div class="modal-content">

      <form role="form" method="post" enctype="multipart/form-data">

          <div class="modal-header" style="background:#3c8dbc; color:white;">
            <button type="button" class="close" data-dismiss="modal">&times;</button>
            <h4 class="modal-title">Agregar Usuario</h4>
          </div>

          <div class="modal-body">

              <div class="box-body">

                <!--ENTRADA PARA EL NOMBRE-->

                  <div class="form-group">
                    <div class="input-group">
                      <span class="input-group-addon"><i class="fa fa-user"></i></span>
                      <input type="text" class="form-control input-lg" name="nuevoNombre" placeholder="Ingresar nombre" required>
                    </div>
                  </div>

                  <!--ENTRADA PARA EL USUARIO-->

                  <div class="form-group">
                    <div class="input-group">
                      <span class="input-group-addon"><i class="fa fa-key"></i></span>
                      <input type="text" class="form-control input-lg" name="nuevoUsuario" placeholder="Ingresar usuario" required>
                    </div>
                  </div>

                  <!--ENTRADA PARA LA CONTRASEÑA-->

                  <div class="form-group">
                    <div class="input-group">
                      <span class="input-group-addon"><i class="fa fa-lock"></i></span>
                      <input type="password" class="form-control input-lg" name="nuevoPassword" placeholder="Ingresar contraseña" required>
                    </div>
                  </div>

                  <!--ENTRADA PARA SELECCIONAR PERFIL-->

                  <div class="form-group">
                    <div class="input-group">
                      <span class="input-group-addon"><i class="fa fa-users"></i></span>
                      <select class="form-control input-lg" name="nuevoPerfil" required>
                        <option value="">Seleccionar perfil</option>
                        <option value="Administrador">Administrador</option>
                        <option value="Especial">Especial</option>
                        <option value="Vendedor">Vendedor</option>
                      </select>
                    </div>
                  </div>

                  <!--ENTRADA PARA SUBIR FOTO-->

                  <div class="form-group">

                    <div class="panel">SUBIR FOTO</div>

                    <input type="file" class="nuevaFoto" name="nuevaFoto" accept="image/*">

                    <p class="help-block">Peso máximo de la foto 2 MB</p>

                    <img src="vistas/img/usuarios/default/anonymous.png" class="img-thumbnail previewFoto" width="100px">

                  </div>

              </div>

          </div>

          <div class="modal-footer">

            <button type="button" class="btn btn-default pull-left" data-dismiss="modal">Salir</button>

            <button type="submit" class="btn btn-primary">Guardar usuario</button>

          </div>

      </form>

    </div>
  </div>
</div>

<!--=====================================
MODAL EDITAR USUARIO
======================================-->

<div class="modal fade" id="modalEditarUsuario">
  <div class="modal-dialog">
    <div class="modal-content">

      <form role="form" method="post" enctype="multipart/form-data">

          <div class="modal-header" style="background:#f39c12; color:white;">
            <button type="button" class="close" data-dismiss="modal">&times;</button>
            <h4 class="modal-title">Editar Usuario</h4>
          </div>

          <div class="modal-body">

              <div class="box-body">

                <!--ENTRADA PARA EL NOMBRE-->

                  <div class="form-group">
                    <div class="input-group">
                      <span class="input-group-addon"><i class="fa fa-user"></i></span>
                      <input type="text" class="form-control input-lg" name="editarNombre" id="editarNombre" placeholder="Nombre" required>
                    </div>
                  </div>

                  <!--ENTRADA PARA EL USUARIO (solo lectura)-->

                  <div class="form-group">
                    <div class="input-group">
                      <span class="input-group-addon"><i class="fa fa-key"></i></span>
                      <input type="text" class="form-control input-lg" name="editarUsuario" id="editarUsuario" readonly>
                    </div>
                  </div>

                  <!--ENTRADA PARA LA CONTRASEÑA-->

                  <div class="form-group">
                    <div class="input-group">
                      <span class="input-group-addon"><i class="fa fa-lock"></i></span>
                      <input type="password" class="form-control input-lg" name="editarPassword" placeholder="Nueva contraseña (dejar vacío para no cambiar)">
                    </div>
                  </div>

                  <input type="hidden" name="passwordActual" id="passwordActual">

                  <!--ENTRADA PARA SELECCIONAR PERFIL-->

                  <div class="form-group">
                    <div class="input-group">
                      <span class="input-group-addon"><i class="fa fa-users"></i></span>
                      <select class="form-control input-lg" name="editarPerfil" id="editarPerfil" required>
                        <option value="Administrador">Administrador</option>
                        <option value="Especial">Especial</option>
                        <option value="Vendedor">Vendedor</option>
                      </select>
                    </div>
                  </div>

                  <!--ENTRADA PARA SUBIR FOTO-->

                  <div class="form-group">

                    <div class="panel">CAMBIAR FOTO</div>

                    <input type="file" class="editarFoto" name="editarFoto" accept="image/*">

                    <p class="help-block">Peso máximo de la foto 2 MB</p>

                    <img src="vistas/img/usuarios/default/anonymous.png" class="img-thumbnail previewFotoEditar" width="100px" id="fotoActualImg">

                  </div>

                  <input type="hidden" name="fotoActual" id="fotoActual">

              </div>

          </div>

          <div class="modal-footer">

            <button type="button" class="btn btn-default pull-left" data-dismiss="modal">Salir</button>

            <button type="submit" class="btn btn-warning">Editar usuario</button>

          </div>

      </form>

    </div>
  </div>
</div>

<?php

$crearUsuario = new ControladorUsuarios();
$crearUsuario -> ctrCrearUsuario();

$editarUsuario = new ControladorUsuarios();
$editarUsuario -> ctrEditarUsuario();

$borrarUsuario = new ControladorUsuarios();
$borrarUsuario -> ctrBorrarUsuario();

?>

<script>

/*=============================================
EDITAR USUARIO - Cargar datos en modal
=============================================*/

$(".tablas").on("click", ".btnEditarUsuario", function(){

  var idUsuario = $(this).attr("idUsuario");

  var datos = new FormData();
  datos.append("idUsuario", idUsuario);

  $.ajax({

    url: "ajax/usuarios.ajax.php",
    method: "POST",
    data: datos,
    cache: false,
    contentType: false,
    processData: false,
    dataType: "json",
    success: function(respuesta){

      $("#editarNombre").val(respuesta["nombre"]);
      $("#editarUsuario").val(respuesta["usuario"]);
      $("#passwordActual").val(respuesta["password"]);
      $("#editarPerfil").val(respuesta["perfil"]);
      $("#fotoActual").val(respuesta["foto"]);

      if(respuesta["foto"] != ""){
        $("#fotoActualImg").attr("src", respuesta["foto"]);
      }else{
        $("#fotoActualImg").attr("src", "vistas/img/usuarios/default/anonymous.png");
      }

    }

  });

});

/*=============================================
ELIMINAR USUARIO
=============================================*/

$(".tablas").on("click", ".btnEliminarUsuario", function(){

  var idUsuario = $(this).attr("idUsuario");
  var usuario = $(this).attr("usuario");

  swal({

    title: "¿Está seguro de eliminar el usuario?",
    text: "¡Si no lo está, puede cancelar la acción!",
    type: "warning",
    showCancelButton: true,
    confirmButtonColor: "#DD6B55",
    cancelButtonText: "Cancelar",
    confirmButtonText: "¡Sí, eliminar usuario!",
    closeOnConfirm: false

  }).then(function(result){

    if(result.value){

      window.location = "usuarios&idUsuario="+idUsuario;

    }

  });

});

/*=============================================
PREVIEW DE IMAGEN AL SELECCIONAR
=============================================*/

$(".nuevaFoto").change(function(){

  var imagen = this.files[0];

  if(imagen["type"] != "image/jpeg" && imagen["type"] != "image/png"){

    $(".nuevaFoto").val("");
    swal({type: "error", title: "¡Error!", text: "¡La imagen debe ser en formato JPG o PNG!"});

  }else if(imagen["size"] > 2000000){

    $(".nuevaFoto").val("");
    swal({type: "error", title: "¡Error!", text: "¡La imagen no debe pesar más de 2 MB!"});

  }else{

    var datosImagen = new FileReader();
    datosImagen.readAsDataURL(imagen);

    $(datosImagen).on("load", function(event){

      var rutaImagen = event.target.result;
      $(".previewFoto").attr("src", rutaImagen);

    });

  }

});

</script>