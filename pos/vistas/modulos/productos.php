<?php

$productos = ControladorProductos::ctrMostrarProductos(null, null);
$categorias = ControladorCategorias::ctrMostrarCategorias(null, null);

?>

<div class="content-wrapper">

  <section class="content-header">

    <h1>
      Administrar productos
    </h1>

    <ol class="breadcrumb">
      <li><a href="inicio"><i class="fa fa-dashboard"></i> Inicio</a></li>
      <li class="active">Productos</li>
    </ol>

  </section>

  <section class="content">

    <div class="box">

      <div class="box-header with-border">

        <button class="btn btn-primary" data-toggle="modal" data-target="#modalAgregarProducto">
          <i class="fa fa-plus"></i> Agregar Producto
        </button>

      </div>

      <div class="box-body">

        <table class="table table-bordered table-striped dt-responsive tablas">

          <thead>
            <tr>
              <th style="width:10px;">#</th>
              <th>Imagen</th>
              <th>Código</th>
              <th>Descripción</th>
              <th>Categoría</th>
              <th>P. Compra</th>
              <th>P. Venta</th>
              <th>Stock</th>
              <th>Estado</th>
              <th>Acciones</th>
            </tr>
          </thead>

          <tbody>

            <?php

            $contador = 1;

            foreach($productos as $key => $valor){

              echo '<tr>

                <td>'.$contador.'</td>
                <td><img src="'.(!empty($valor["imagen"]) ? $valor["imagen"] : 'vistas/img/productos/default/anonymous.png').'" class="img-thumbnail" width="40px"></td>
                <td>'.$valor["codigo"].'</td>
                <td>'.$valor["descripcion"].'</td>
                <td><span class="label label-info">'.$valor["categoria"].'</span></td>
                <td>$'.number_format($valor["precio_compra"], 2).'</td>
                <td>$'.number_format($valor["precio_venta"], 2).'</td>
                <td>';

              if($valor["stock"] <= $valor["stock_minimo"]){
                echo '<span class="badge bg-red">'.$valor["stock"].'</span>';
              }else{
                echo '<span class="badge bg-green">'.$valor["stock"].'</span>';
              }

              echo '</td>
                <td>';

              if($valor["estado"] == 1){
                echo '<button class="btn btn-success btn-xs">Activado</button>';
              }else{
                echo '<button class="btn btn-danger btn-xs">Inactivo</button>';
              }

              echo '</td>
                <td>

                  <div class="btn-group">

                    <button class="btn btn-warning btnEditarProducto" data-toggle="modal" data-target="#modalEditarProducto" idProducto="'.$valor["id"].'"><i class="fa fa-pencil"></i></button>

                    <a class="btn btn-danger btnEliminarProducto" idProducto="'.$valor["id"].'"><i class="fa fa-times"></i></a>

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
MODAL AGREGAR PRODUCTO
======================================-->

<div class="modal fade" id="modalAgregarProducto">
  <div class="modal-dialog">
    <div class="modal-content">

      <form role="form" method="post" enctype="multipart/form-data">

          <div class="modal-header" style="background:#3c8dbc; color:white;">
            <button type="button" class="close" data-dismiss="modal">&times;</button>
            <h4 class="modal-title">Agregar Producto</h4>
          </div>

          <div class="modal-body">

              <div class="box-body">

                <!-- CODIGO -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-code"></i></span>
                    <input type="text" class="form-control input-lg" name="nuevoCodigo" placeholder="Código del producto" required>
                  </div>
                </div>

                <!-- CODIGO DE BARRAS -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-barcode"></i></span>
                    <input type="text" class="form-control input-lg" name="nuevoCodigoBarras" placeholder="Código de barras (opcional)">
                  </div>
                </div>

                <!-- DESCRIPCION -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-product-hunt"></i></span>
                    <input type="text" class="form-control input-lg" name="nuevaDescripcion" placeholder="Descripción del producto" required>
                  </div>
                </div>

                <!-- CATEGORIA -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-th"></i></span>
                    <select class="form-control input-lg" name="nuevaCategoria" required>
                      <option value="">Seleccionar categoría</option>
                      <?php

                      foreach($categorias as $key => $valor){

                        if($valor["estado"] == 1){
                          echo '<option value="'.$valor["id"].'">'.$valor["nombre"].'</option>';
                        }

                      }

                      ?>
                    </select>
                  </div>
                </div>

                <!-- PRECIOS -->
                <div class="row">

                  <div class="col-xs-6">
                    <div class="form-group">
                      <div class="input-group">
                        <span class="input-group-addon"><i class="fa fa-arrow-down"></i></span>
                        <input type="number" class="form-control input-lg" name="nuevoPrecioCompra" step="0.01" min="0" placeholder="Precio compra" required>
                      </div>
                    </div>
                  </div>

                  <div class="col-xs-6">
                    <div class="form-group">
                      <div class="input-group">
                        <span class="input-group-addon"><i class="fa fa-arrow-up"></i></span>
                        <input type="number" class="form-control input-lg" name="nuevoPrecioVenta" step="0.01" min="0" placeholder="Precio venta" required>
                      </div>
                    </div>
                  </div>

                </div>

                <!-- STOCK -->
                <div class="row">

                  <div class="col-xs-6">
                    <div class="form-group">
                      <div class="input-group">
                        <span class="input-group-addon"><i class="fa fa-cubes"></i></span>
                        <input type="number" class="form-control input-lg" name="nuevoStock" min="0" value="0" placeholder="Stock" required>
                      </div>
                    </div>
                  </div>

                  <div class="col-xs-6">
                    <div class="form-group">
                      <div class="input-group">
                        <span class="input-group-addon"><i class="fa fa-exclamation-triangle"></i></span>
                        <input type="number" class="form-control input-lg" name="nuevoStockMinimo" min="0" value="5" placeholder="Stock mínimo" required>
                      </div>
                    </div>
                  </div>

                </div>

                <!-- IMAGEN -->
                <div class="form-group">
                  <div class="panel">SUBIR IMAGEN</div>
                  <input type="file" class="nuevaImagen" name="nuevaImagen" accept="image/*">
                  <p class="help-block">Peso máximo 2 MB. Formato JPG o PNG.</p>
                  <img src="vistas/img/productos/default/anonymous.png" class="img-thumbnail previewImagen" width="100px">
                </div>

              </div>

          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-default pull-left" data-dismiss="modal">Salir</button>
            <button type="submit" class="btn btn-primary">Guardar producto</button>
          </div>

      </form>

    </div>
  </div>
</div>

<!--=====================================
MODAL EDITAR PRODUCTO
======================================-->

<div class="modal fade" id="modalEditarProducto">
  <div class="modal-dialog">
    <div class="modal-content">

      <form role="form" method="post" enctype="multipart/form-data">

          <div class="modal-header" style="background:#f39c12; color:white;">
            <button type="button" class="close" data-dismiss="modal">&times;</button>
            <h4 class="modal-title">Editar Producto</h4>
          </div>

          <div class="modal-body">

              <div class="box-body">

                <input type="hidden" name="idProductoEditar" id="idProductoEditar">
                <input type="hidden" name="imagenActual" id="imagenActual">

                <!-- CODIGO -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-code"></i></span>
                    <input type="text" class="form-control input-lg" name="editarCodigo" id="editarCodigo" placeholder="Código" required>
                  </div>
                </div>

                <!-- CODIGO DE BARRAS -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-barcode"></i></span>
                    <input type="text" class="form-control input-lg" name="editarCodigoBarras" id="editarCodigoBarras" placeholder="Código de barras">
                  </div>
                </div>

                <!-- DESCRIPCION -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-product-hunt"></i></span>
                    <input type="text" class="form-control input-lg" name="editarDescripcion" id="editarDescripcion" placeholder="Descripción" required>
                  </div>
                </div>

                <!-- CATEGORIA -->
                <div class="form-group">
                  <div class="input-group">
                    <span class="input-group-addon"><i class="fa fa-th"></i></span>
                    <select class="form-control input-lg" name="editarCategoria" id="editarCategoria" required>
                      <?php

                      foreach($categorias as $key => $valor){

                        if($valor["estado"] == 1){
                          echo '<option value="'.$valor["id"].'">'.$valor["nombre"].'</option>';
                        }

                      }

                      ?>
                    </select>
                  </div>
                </div>

                <!-- PRECIOS -->
                <div class="row">

                  <div class="col-xs-6">
                    <div class="form-group">
                      <div class="input-group">
                        <span class="input-group-addon"><i class="fa fa-arrow-down"></i></span>
                        <input type="number" class="form-control input-lg" name="editarPrecioCompra" id="editarPrecioCompra" step="0.01" min="0" required>
                      </div>
                    </div>
                  </div>

                  <div class="col-xs-6">
                    <div class="form-group">
                      <div class="input-group">
                        <span class="input-group-addon"><i class="fa fa-arrow-up"></i></span>
                        <input type="number" class="form-control input-lg" name="editarPrecioVenta" id="editarPrecioVenta" step="0.01" min="0" required>
                      </div>
                    </div>
                  </div>

                </div>

                <!-- STOCK -->
                <div class="row">

                  <div class="col-xs-6">
                    <div class="form-group">
                      <div class="input-group">
                        <span class="input-group-addon"><i class="fa fa-cubes"></i></span>
                        <input type="number" class="form-control input-lg" name="editarStock" id="editarStock" min="0" required>
                      </div>
                    </div>
                  </div>

                  <div class="col-xs-6">
                    <div class="form-group">
                      <div class="input-group">
                        <span class="input-group-addon"><i class="fa fa-exclamation-triangle"></i></span>
                        <input type="number" class="form-control input-lg" name="editarStockMinimo" id="editarStockMinimo" min="0" required>
                      </div>
                    </div>
                  </div>

                </div>

                <!-- IMAGEN -->
                <div class="form-group">
                  <div class="panel">CAMBIAR IMAGEN</div>
                  <input type="file" class="editarImagen" name="editarImagen" accept="image/*">
                  <p class="help-block">Peso máximo 2 MB. Formato JPG o PNG.</p>
                  <img src="vistas/img/productos/default/anonymous.png" class="img-thumbnail previewImagenEditar" width="100px" id="imagenActualImg">
                </div>

              </div>

          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-default pull-left" data-dismiss="modal">Salir</button>
            <button type="submit" class="btn btn-warning">Editar producto</button>
          </div>

      </form>

    </div>
  </div>
</div>

<?php

$crearProducto = new ControladorProductos();
$crearProducto -> ctrCrearProducto();

$editarProducto = new ControladorProductos();
$editarProducto -> ctrEditarProducto();

$borrarProducto = new ControladorProductos();
$borrarProducto -> ctrBorrarProducto();

?>

<script>

/*=============================================
EDITAR PRODUCTO - Cargar datos en modal via AJAX
=============================================*/

$(".tablas").on("click", ".btnEditarProducto", function(){

  var idProducto = $(this).attr("idProducto");

  var datos = new FormData();
  datos.append("idProducto", idProducto);

  $.ajax({

    url: "ajax/productos.ajax.php",
    method: "POST",
    data: datos,
    cache: false,
    contentType: false,
    processData: false,
    dataType: "json",
    success: function(respuesta){

      $("#idProductoEditar").val(respuesta["id"]);
      $("#editarCodigo").val(respuesta["codigo"]);
      $("#editarCodigoBarras").val(respuesta["codigo_barras"]);
      $("#editarDescripcion").val(respuesta["descripcion"]);
      $("#editarCategoria").val(respuesta["id_categoria"]);
      $("#editarPrecioCompra").val(respuesta["precio_compra"]);
      $("#editarPrecioVenta").val(respuesta["precio_venta"]);
      $("#editarStock").val(respuesta["stock"]);
      $("#editarStockMinimo").val(respuesta["stock_minimo"]);
      $("#imagenActual").val(respuesta["imagen"]);

      if(respuesta["imagen"] != "" && respuesta["imagen"] != null){
        $("#imagenActualImg").attr("src", respuesta["imagen"]);
      }else{
        $("#imagenActualImg").attr("src", "vistas/img/productos/default/anonymous.png");
      }

    }

  });

});

/*=============================================
ELIMINAR PRODUCTO
=============================================*/

$(".tablas").on("click", ".btnEliminarProducto", function(){

  var idProducto = $(this).attr("idProducto");

  swal({

    title: "¿Está seguro de eliminar el producto?",
    text: "¡Si no lo está, puede cancelar la acción!",
    type: "warning",
    showCancelButton: true,
    confirmButtonColor: "#DD6B55",
    cancelButtonText: "Cancelar",
    confirmButtonText: "¡Sí, eliminar!",
    closeOnConfirm: false

  }).then(function(result){

    if(result.value){

      window.location = "productos&idProducto="+idProducto;

    }

  });

});

/*=============================================
PREVIEW DE IMAGEN AL SELECCIONAR
=============================================*/

$(".nuevaImagen").change(function(){

  var imagen = this.files[0];

  if(imagen["type"] != "image/jpeg" && imagen["type"] != "image/png"){

    $(".nuevaImagen").val("");
    swal({type: "error", title: "¡Error!", text: "¡La imagen debe ser en formato JPG o PNG!"});

  }else if(imagen["size"] > 2000000){

    $(".nuevaImagen").val("");
    swal({type: "error", title: "¡Error!", text: "¡La imagen no debe pesar más de 2 MB!"});

  }else{

    var datosImagen = new FileReader();
    datosImagen.readAsDataURL(imagen);

    $(datosImagen).on("load", function(event){
      var rutaImagen = event.target.result;
      $(".previewImagen").attr("src", rutaImagen);
    });

  }

});

</script>
