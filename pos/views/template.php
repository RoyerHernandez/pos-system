<?php

session_start();

?>

<!DOCTYPE html>
<html>
<head>

  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">

  <title>Callejon Bar - POS</title>

  <!-- Tell the browser to be responsive to screen width -->
  <meta content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" name="viewport">

  <link rel="icon" href="views/img/plantilla/favicon.svg" type="image/svg+xml">

   <!--=====================================
  CSS PLUGINS
  ======================================-->

  <!-- Bootstrap 3.3.7 -->
  <link rel="stylesheet" href="views/bower_components/bootstrap/dist/css/bootstrap.min.css">

  <!-- Font Awesome -->
  <link rel="stylesheet" href="views/bower_components/font-awesome/css/font-awesome.min.css">

  <!-- Ionicons -->
  <link rel="stylesheet" href="views/bower_components/Ionicons/css/ionicons.min.css">

  <!-- Theme style -->
  <link rel="stylesheet" href="views/dist/css/AdminLTE.css">

  <!-- AdminLTE Skins -->
  <link rel="stylesheet" href="views/dist/css/skins/_all-skins.min.css">

  <!-- Callejon Bar Branding -->
  <link rel="stylesheet" href="views/css/callejon-brand.css">

  <!-- Google Font -->
  <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Source+Sans+Pro:300,400,600,700,300italic,400italic,600italic">

   <!-- DataTables -->
  <link rel="stylesheet" href="views/bower_components/datatables.net-bs/css/dataTables.bootstrap.min.css">

  <!--=====================================
  JAVASCRIPT PLUGINS
  ======================================-->

  <!-- jQuery 3 -->
  <script src="views/bower_components/jquery/dist/jquery.min.js"></script>

  <!-- Bootstrap 3.3.7 -->
  <script src="views/bower_components/bootstrap/dist/js/bootstrap.min.js"></script>

  <!-- FastClick -->
  <script src="views/bower_components/fastclick/lib/fastclick.js"></script>

  <!-- AdminLTE App -->
  <script src="views/dist/js/adminlte.min.js"></script>

  <!-- DataTables -->
  <script src="views/bower_components/datatables.net/js/jquery.dataTables.min.js"></script>
  <script src="views/bower_components/datatables.net-bs/js/dataTables.bootstrap.min.js"></script>

  <!-- Chart.js v2.9.4 -->
  <script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.9.4/Chart.min.js"></script>

  <!-- SweetAlert 2 -->
  <script src="views/plugins/sweetalert2/sweetalert2.all.js"></script>
  <!-- By default SweetAlert2 doesn't support IE. To enable IE 11 support, include Promise polyfill:-->
  <script src="https://cdnjs.cloudflare.com/ajax/libs/core-js/2.4.1/core.js"></script>

</head>

<!--=====================================
DOCUMENT BODY
======================================-->

<body class="hold-transition skin-black sidebar-collapse sidebar-mini login-page">

  <?php

  if(isset($_SESSION["loggedIn"]) && $_SESSION["loggedIn"] == "ok"){

   echo '<div class="wrapper">';

    /*=============================================
    HEADER
    =============================================*/

    include "modules/header.php";

    /*=============================================
    MENU
    =============================================*/

    include "modules/menu.php";

    /*=============================================
    CONTENT
    =============================================*/

    if(isset($_GET["ruta"])){

      /*=============================================
      ROLE-BASED ACCESS CONTROL
      =============================================*/

      $role = isset($_SESSION["role"]) ? $_SESSION["role"] : "Vendedor";

      // Routes accessible by all roles
      $allRoles = array("inicio", "crear-venta", "caja", "salir");

      // Routes for Admin and Especial
      $adminEspecial = array("productos", "clientes", "reportes", "ventas", "inventario");

      // Routes for Admin only
      $adminOnly = array("usuarios", "categorias");

      $route = $_GET["ruta"];
      $allowed = false;

      if(in_array($route, $allRoles)){

        $allowed = true;

      }else if(in_array($route, $adminEspecial)){

        if($role == "Administrador" || $role == "Especial"){
          $allowed = true;
        }

      }else if(in_array($route, $adminOnly)){

        if($role == "Administrador"){
          $allowed = true;
        }

      }

      if($allowed){

        include "modules/".$route.".php";

      }else{

        include "modules/403.php";

      }

    }else{

      include "modules/inicio.php";

    }

    /*=============================================
    FOOTER
    =============================================*/

    include "modules/footer.php";

    echo '</div>';

  }else{

    include "modules/login.php";

  }

  ?>


<script src="views/js/plantilla.js"></script>
<script src="views/js/sales.js"></script>
<script src="views/js/sales-admin.js"></script>
<script src="views/js/dashboard.js"></script>
<script src="views/js/reports.js"></script>
<script src="views/js/inventory.js"></script>

</body>
</html>