/*=============================================
SESSION TIMEOUT — 2 minutes of inactivity
Resets on any user interaction.
Shows a warning 30 s before logout.
=============================================*/

if (!$("body").hasClass("login-page")) {

    var SESSION_MS  = 120000; // 2 minutes
    var WARNING_MS  =  30000; // warn at 30 s remaining
    var lastActivity = Date.now();
    var warningShown = false;

    // Reset inactivity clock on any user interaction
    $(document).on("mousemove click keypress scroll touchstart", function() {
        lastActivity = Date.now();
        if (warningShown) {
            warningShown = false;
            swal.close();
        }
    });

    setInterval(function() {

        var remaining = SESSION_MS - (Date.now() - lastActivity);

        if (remaining <= 0) {

            window.location = "salir";

        } else if (remaining <= WARNING_MS && !warningShown) {

            warningShown = true;

            swal({
                type: "warning",
                title: "¡Sesión por expirar!",
                text: "Su sesión se cerrará en 30 segundos por inactividad.",
                showConfirmButton: true,
                confirmButtonText: "Continuar sesión",
                allowOutsideClick: false
            }).then(function(result) {
                if (result.value) {
                    lastActivity = Date.now();
                    warningShown = false;
                }
            });

        }

    }, 1000);

}

/*=============================================
SideBar Menu
=============================================*/

$('.sidebar-menu').tree()

/*=============================================
Data Table
=============================================*/

$(".tablas").DataTable({

	"pagingType": "simple_numbers",

	"language": {

		"sProcessing":     "Procesando...",
		"sLengthMenu":     "Mostrar _MENU_ registros",
		"sZeroRecords":    "No se encontraron resultados",
		"sEmptyTable":     "Ningún dato disponible en esta tabla",
		"sInfo":           "Mostrando registros del _START_ al _END_ de un total de _TOTAL_",
		"sInfoEmpty":      "Mostrando registros del 0 al 0 de un total de 0",
		"sInfoFiltered":   "(filtrado de un total de _MAX_ registros)",
		"sInfoPostFix":    "",
		"sSearch":         "Buscar:",
		"sUrl":            "",
		"sInfoThousands":  ",",
		"sLoadingRecords": "Cargando...",
		"oPaginate": {
		"sFirst":    "Primero",
		"sLast":     "Último",
		"sNext":     "Siguiente",
		"sPrevious": "Anterior"
		},
		"oAria": {
			"sSortAscending":  ": Activar para ordenar la columna de manera ascendente",
			"sSortDescending": ": Activar para ordenar la columna de manera descendente"
		}

	}

});