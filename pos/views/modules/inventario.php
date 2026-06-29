<div class="content-wrapper">

  <section class="content-header">
    <h1>
      Administracion de Inventario
      <small>Entradas y salidas de productos</small>
    </h1>
    <ol class="breadcrumb">
      <li><a href="inicio"><i class="fa fa-dashboard"></i> Inicio</a></li>
      <li class="active">Inventario</li>
    </ol>
  </section>

  <section class="content">

    <!-- Action Bar -->
    <div class="row" style="margin-bottom:15px;">
      <div class="col-md-12">
        <button class="btn btn-primary" data-toggle="modal" data-target="#modalNewMovement">
          <i class="fa fa-plus"></i> Nuevo Movimiento
        </button>
      </div>
    </div>

    <!-- Date Range Filter -->
    <div class="box box-info">
      <div class="box-header with-border">
        <h3 class="box-title"><i class="fa fa-calendar"></i> Filtrar por Fecha</h3>
      </div>
      <div class="box-body">
        <form class="form-inline" id="formFilterMovements">
          <div class="form-group">
            <label>Desde:</label>
            <input type="date" class="form-control" id="movStartDate" value="<?php echo date('Y-m-01'); ?>">
          </div>
          <div class="form-group" style="margin-left:15px;">
            <label>Hasta:</label>
            <input type="date" class="form-control" id="movEndDate" value="<?php echo date('Y-m-d'); ?>">
          </div>
          <button type="submit" class="btn btn-primary" style="margin-left:15px;">
            <i class="fa fa-search"></i> Filtrar
          </button>
          <button type="button" class="btn btn-default" id="btnClearFilter" style="margin-left:5px;">
            <i class="fa fa-times"></i> Limpiar
          </button>
        </form>
      </div>
    </div>

    <!-- Movements Table -->
    <div class="box box-primary">
      <div class="box-header with-border">
        <h3 class="box-title"><i class="fa fa-exchange"></i> Historial de Movimientos</h3>
      </div>
      <div class="box-body">
        <table class="table table-bordered table-striped dt-responsive" id="tableMovements" width="100%">
          <thead>
            <tr>
              <th>#</th>
              <th>Fecha</th>
              <th>Producto</th>
              <th>Tipo</th>
              <th>Motivo</th>
              <th>Cantidad</th>
              <th>Usuario</th>
              <th>Observaciones</th>
            </tr>
          </thead>
          <tbody id="movementsBody">
          </tbody>
        </table>
      </div>
    </div>

  </section>

</div>

<!--=============================================
MODAL: NEW MOVEMENT
=============================================-->

<div class="modal fade" id="modalNewMovement" tabindex="-1" role="dialog">
  <div class="modal-dialog" role="document">
    <div class="modal-content">

      <form id="formNewMovement">

        <div class="modal-header" style="background:#3c8dbc; color:white;">
          <button type="button" class="close" data-dismiss="modal" style="color:white;">&times;</button>
          <h4 class="modal-title"><i class="fa fa-exchange"></i> Nuevo Movimiento de Inventario</h4>
        </div>

        <div class="modal-body">

          <!-- Product Search -->
          <div class="form-group">
            <label>Producto</label>
            <div style="position:relative;">
              <input type="text" class="form-control" id="searchMovProduct" placeholder="Buscar producto por codigo o nombre..." autocomplete="off">
              <input type="hidden" id="movProductId">
              <div id="movProductResults" class="list-group" style="position:absolute; z-index:1000; width:100%; display:none; max-height:200px; overflow-y:auto;">
              </div>
            </div>
            <p id="movProductInfo" class="help-block" style="margin-top:5px;"></p>
          </div>

          <!-- Type -->
          <div class="form-group">
            <label>Tipo de Movimiento</label>
            <select class="form-control" id="movType" required>
              <option value="">-- Seleccionar --</option>
              <option value="entrada">Entrada</option>
              <option value="salida">Salida</option>
            </select>
          </div>

          <!-- Motivo -->
          <div class="form-group">
            <label>Motivo</label>
            <select class="form-control" id="movMotivo" required disabled>
              <option value="">-- Seleccione tipo primero --</option>
            </select>
          </div>

          <!-- Cantidad -->
          <div class="form-group">
            <label>Cantidad</label>
            <input type="number" class="form-control" id="movCantidad" min="1" value="1" required>
          </div>

          <!-- Observaciones -->
          <div class="form-group">
            <label>Observaciones <small class="text-muted">(opcional)</small></label>
            <textarea class="form-control" id="movObservaciones" rows="2" placeholder="Notas adicionales..."></textarea>
          </div>

        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-default" data-dismiss="modal">Cancelar</button>
          <button type="submit" class="btn btn-primary" id="btnSaveMovement">
            <i class="fa fa-save"></i> Registrar Movimiento
          </button>
        </div>

      </form>

    </div>
  </div>
</div>
