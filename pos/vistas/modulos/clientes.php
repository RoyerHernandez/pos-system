<?php

$clientes = ControladorClientes::ctrMostrarClientes(null, null);

?>

<div class="content-wrapper">

  <section class="content-header">

    <h1>
      Administrar clientes
    </h1>

    <ol class="breadcrumb">
      <li><a href="inicio"><i class="fa fa-dashboard"></i> Inicio</a></li>
      <li class="active">Clientes</li>
    </ol>

  </section>

  <section class="content">

    <div class="box">

      <div class="box-header with-border">

        <button class="btn btn-primary" data-toggle="modal" data-target="#modalAgregarCliente">
          <i class="fa fa-plus"></i> Agregar Cliente
        </button>

      </div>

      <div class="box-body">

        <table class="table table-bordered table-striped dt-responsive tablas">

          <thead>
            <tr>
              <th style="width:10px;">#</th>
              <th>Nombre</th>
              <th>Documento</th>
              <th>Teléfono</th>
              <th>Email</th>
              <th>Total Compras</th>
              <th>Estado</th>
              <th>Acciones</th>
            </tr>
          </thead>

          <tbody>

            <?php

            $contador = 1;

            foreach($clientes as $key => $valor){

              echo '<tr>

                <td>'.$contador.'</td>
                <td>'.$valor["nombre"].'</td>
                <td>'.$valor["documento"].'</td>
                <td>'.$valor["telefono"].'</td>
                <td>'.$valor["email"].'</td>
                <td>$'.number_format($valor["total_compras"], 2).'</td>
                <td>';

              if($valor["estado"] == 1){
                echo '<button class="btn btn-success btn-xs">Activado</button>';
              }else{
                echo '<button class="btn btn-danger btn-xs">Inactivo</button>';
              }

              echo '</td>
                <td>

                  <div class="btn-group">';

              // No permitir eliminar "Publico General"
              if($valor["id"] != 1){

                echo '<button class="btn btn-warning btnEditarCliente" data-toggle="modal" data-target="#modalEditarCliente" idCliente="'.$valor["id"].'"><i class="fa fa-pencil"></i></button>

                    <a class="btn btn-danger btnEliminarCliente" idCliente="'.$valor["id"].'"><i class="fa fa-times"></i></a>';

              }else{

                echo '<button class="btn btn-default btn-xs" disabled><i class="fa fa-lock"></i> Cliente por defecto</button>';

              }

              echo '</div>

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
MODAL AGREGAR CLIENTE
======================================-->

<div class="modal fade" id="modalAgregarCliente">
  <div class="modal-dialog">
    <div class="modal-content">

      <form role="form" method="post">

          <div class="modal-header" style="background:#3c8dbc; color:white;">
            <button type="button" class="close" data-dismiss="modal">&times;</button>
            <h4 class="modal-title">Agregar Cliente</h4>
          </div>

          <div class="modal-body">

              <div class="box-body">

                <!-- NOMBRE -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-user"></i></span>
                    <input type="text" class="form-control input-lg" name="nuevoCliente" placeholder="Nombre del cliente" required>
                  </div>
                </div>

                <!-- DOCUMENTO -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-id-card"></i></span>
                    <input type="text" class="form-control input-lg" name="nuevoDocumento" placeholder="Documento / NIT (opcional)">
                  </div>
                </div>

                <!-- TELEFONO -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-phone"></i></span>
                    <input type="text" class="form-control input-lg" name="nuevoTelefono" placeholder="Teléfono (opcional)">
                  </div>
                </div>

                <!-- EMAIL -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-envelope"></i></span>
                    <input type="email" class="form-control input-lg" name="nuevoEmail" placeholder="Email (opcional)">
                  </div>
                </div>

                <!-- DIRECCION -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-map-marker"></i></span>
                    <input type="text" class="form-control input-lg" name="nuevaDireccion" placeholder="Dirección (opcional)">
                  </div>
                </div>

                <!-- FECHA DE NACIMIENTO -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-calendar"></i></span>
                    <input type="date" class="form-control input-lg" name="nuevaFechaNacimiento">
                  </div>
                </div>

              </div>

          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-default pull-left" data-dismiss="modal">Salir</button>
            <button type="submit" class="btn btn-primary">Guardar cliente</button>
          </div>

      </form>

    </div>
  </div>
</div>

<!--=====================================
MODAL EDITAR CLIENTE
======================================-->

<div class="modal fade" id="modalEditarCliente">
  <div class="modal-dialog">
    <div class="modal-content">

      <form role="form" method="post">

          <div class="modal-header" style="background:#f39c12; color:white;">
            <button type="button" class="close" data-dismiss="modal">&times;</button>
            <h4 class="modal-title">Editar Cliente</h4>
          </div>

          <div class="modal-body">

              <div class="box-body">

                <input type="hidden" name="idClienteEditar" id="idClienteEditar">

                <!-- NOMBRE -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-user"></i></span>
                    <input type="text" class="form-control input-lg" name="editarCliente" id="editarCliente" placeholder="Nombre" required>
                  </div>
                </div>

                <!-- DOCUMENTO -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-id-card"></i></span>
                    <input type="text" class="form-control input-lg" name="editarDocumento" id="editarDocumento" placeholder="Documento / NIT">
                  </div>
                </div>

                <!-- TELEFONO -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-phone"></i></span>
                    <input type="text" class="form-control input-lg" name="editarTelefono" id="editarTelefono" placeholder="Teléfono">
                  </div>
                </div>

                <!-- EMAIL -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-envelope"></i></span>
                    <input type="email" class="form-control input-lg" name="editarEmail" id="editarEmail" placeholder="Email">
                  </div>
                </div>

                <!-- DIRECCION -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-map-marker"></i></span>
                    <input type="text" class="form-control input-lg" name="editarDireccion" id="editarDireccion" placeholder="Dirección">
                  </div>
                </div>

                <!-- FECHA DE NACIMIENTO -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-calendar"></i></span>
                    <input type="date" class="form-control input-lg" name="editarFechaNacimiento" id="editarFechaNacimiento">
                  </div>
                </div>

              </div>

          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-default pull-left" data-dismiss="modal">Salir</button>
            <button type="submit" class="btn btn-warning">Editar cliente</button>
          </div>

      </form>

    </div>
  </div>
</div>

<?php

$crearCliente = new ControladorClientes();
$crearCliente -> ctrCrearCliente();

$editarCliente = new ControladorClientes();
$editarCliente -> ctrEditarCliente();

$borrarCliente = new ControladorClientes();
$borrarCliente -> ctrBorrarCliente();

?>

<script>

/*=============================================
EDITAR CLIENTE - Cargar datos en modal via AJAX
=============================================*/

$(".tablas").on("click", ".btnEditarCliente", function(){

  var idCliente = $(this).attr("idCliente");

  var datos = new FormData();
  datos.append("idCliente", idCliente);

  $.ajax({

    url: "ajax/clientes.ajax.php",
    method: "POST",
    data: datos,
    cache: false,
    contentType: false,
    processData: false,
    dataType: "json",
    success: function(respuesta){

      $("#idClienteEditar").val(respuesta["id"]);
      $("#editarCliente").val(respuesta["nombre"]);
      $("#editarDocumento").val(respuesta["documento"]);
      $("#editarTelefono").val(respuesta["telefono"]);
      $("#editarEmail").val(respuesta["email"]);
      $("#editarDireccion").val(respuesta["direccion"]);
      $("#editarFechaNacimiento").val(respuesta["fecha_nacimiento"]);

    }

  });

});

/*=============================================
ELIMINAR CLIENTE
=============================================*/

$(".tablas").on("click", ".btnEliminarCliente", function(){

  var idCliente = $(this).attr("idCliente");

  swal({

    title: "¿Está seguro de eliminar el cliente?",
    text: "¡Si no lo está, puede cancelar la acción!",
    type: "warning",
    showCancelButton: true,
    confirmButtonColor: "#DD6B55",
    cancelButtonText: "Cancelar",
    confirmButtonText: "¡Sí, eliminar!",
    closeOnConfirm: false

  }).then(function(result){

    if(result.value){

      window.location = "clientes&idCliente="+idCliente;

    }

  });

});

</script>
