<div class="content-wrapper">

  <section class="content-header">
    <h1>
      Reportes de Ventas
      <small>Analisis y exportacion</small>
    </h1>
    <ol class="breadcrumb">
      <li><a href="inicio"><i class="fa fa-dashboard"></i> Inicio</a></li>
      <li class="active">Reportes</li>
    </ol>
  </section>

  <section class="content">

    <!-- Date Range Filter (shared) -->
    <div class="box box-info">
      <div class="box-header with-border">
        <h3 class="box-title"><i class="fa fa-calendar"></i> Periodo del Reporte</h3>
      </div>
      <div class="box-body">
        <form class="form-inline" id="formReportDates">
          <div class="form-group">
            <label>Desde:</label>
            <input type="date" class="form-control" id="reportStartDate" value="<?php echo date('Y-m-01'); ?>">
          </div>
          <div class="form-group" style="margin-left:15px;">
            <label>Hasta:</label>
            <input type="date" class="form-control" id="reportEndDate" value="<?php echo date('Y-m-d'); ?>">
          </div>
          <button type="submit" class="btn btn-primary" style="margin-left:15px;">
            <i class="fa fa-search"></i> Generar Reporte
          </button>
        </form>
      </div>
    </div>

    <!-- Report Tabs -->
    <div class="nav-tabs-custom">

      <ul class="nav nav-tabs">
        <li class="active"><a href="#tabSalesByDate" data-toggle="tab"><i class="fa fa-calendar"></i> Ventas por Fecha</a></li>
        <li><a href="#tabTopProducts" data-toggle="tab"><i class="fa fa-trophy"></i> Top Productos</a></li>
        <li><a href="#tabByPaymentMethod" data-toggle="tab"><i class="fa fa-credit-card"></i> Por Metodo de Pago</a></li>
        <li><a href="#tabCashRegister" data-toggle="tab"><i class="fa fa-desktop"></i> Cierre de Caja</a></li>
      </ul>

      <div class="tab-content">

        <!--=============================================
        TAB 1: SALES BY DATE
        =============================================-->

        <div class="tab-pane active" id="tabSalesByDate">

          <!-- Summary Cards -->
          <div class="row" id="salesSummaryCards" style="display:none; margin-top:15px;">
            <div class="col-md-3">
              <div class="info-box bg-green">
                <span class="info-box-icon"><i class="fa fa-shopping-cart"></i></span>
                <div class="info-box-content">
                  <span class="info-box-text">Ventas Completadas</span>
                  <span class="info-box-number" id="summaryCompleted">0</span>
                  <span class="info-box-text" id="summaryCompletedAmount">$0.00</span>
                </div>
              </div>
            </div>
            <div class="col-md-3">
              <div class="info-box bg-red">
                <span class="info-box-icon"><i class="fa fa-ban"></i></span>
                <div class="info-box-content">
                  <span class="info-box-text">Canceladas</span>
                  <span class="info-box-number" id="summaryCancelled">0</span>
                </div>
              </div>
            </div>
            <div class="col-md-3">
              <div class="info-box bg-aqua">
                <span class="info-box-icon"><i class="fa fa-list"></i></span>
                <div class="info-box-content">
                  <span class="info-box-text">Total Transacciones</span>
                  <span class="info-box-number" id="summaryTotal">0</span>
                </div>
              </div>
            </div>
            <div class="col-md-3">
              <div class="info-box bg-yellow">
                <span class="info-box-icon"><i class="fa fa-dollar"></i></span>
                <div class="info-box-content">
                  <span class="info-box-text">Promedio por Venta</span>
                  <span class="info-box-number" id="summaryAvg">$0.00</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Daily Chart -->
          <div id="dailyChartContainer" style="display:none; margin-bottom:20px;">
            <canvas id="chartDailySales" height="80"></canvas>
          </div>

          <!-- Sales Table -->
          <div id="salesByDateTable" style="margin-top:15px;">
            <div class="text-center text-muted" style="padding:40px;">
              <i class="fa fa-calendar fa-3x"></i>
              <p style="margin-top:10px;">Seleccione un rango de fechas y presione "Generar Reporte"</p>
            </div>
          </div>

        </div>

        <!--=============================================
        TAB 2: TOP PRODUCTS
        =============================================-->

        <div class="tab-pane" id="tabTopProducts">

          <div id="topProductsChart" style="display:none; margin:15px 0;">
            <canvas id="chartTopProductsReport" height="100"></canvas>
          </div>

          <div id="topProductsTable" style="margin-top:15px;">
            <div class="text-center text-muted" style="padding:40px;">
              <i class="fa fa-trophy fa-3x"></i>
              <p style="margin-top:10px;">Seleccione un rango de fechas y presione "Generar Reporte"</p>
            </div>
          </div>

        </div>

        <!--=============================================
        TAB 3: BY PAYMENT METHOD
        =============================================-->

        <div class="tab-pane" id="tabByPaymentMethod">

          <div class="row" style="margin-top:15px;">
            <div class="col-md-5">
              <div id="paymentChartContainer" style="display:none;">
                <canvas id="chartPaymentReport" height="250"></canvas>
              </div>
            </div>
            <div class="col-md-7">
              <div id="paymentMethodTable">
                <div class="text-center text-muted" style="padding:40px;">
                  <i class="fa fa-credit-card fa-3x"></i>
                  <p style="margin-top:10px;">Seleccione un rango de fechas y presione "Generar Reporte"</p>
                </div>
              </div>
            </div>
          </div>

        </div>

        <!--=============================================
        TAB 4: CASH REGISTER
        =============================================-->

        <div class="tab-pane" id="tabCashRegister">

          <div id="cashRegisterTable" style="margin-top:15px;">
            <div class="text-center text-muted" style="padding:40px;">
              <i class="fa fa-desktop fa-3x"></i>
              <p style="margin-top:10px;">Seleccione un rango de fechas y presione "Generar Reporte"</p>
            </div>
          </div>

        </div>

      </div>

    </div>

  </section>

</div>
