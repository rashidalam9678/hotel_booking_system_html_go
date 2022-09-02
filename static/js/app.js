let attention = Prompt();

        // function to validate the forms
        function formValidation() {
          "use strict";

          // Fetch all the forms we want to apply custom Bootstrap validation styles to
          const forms = document.querySelectorAll(".needs-validation");

          // Loop over them and prevent submission
          Array.from(forms).forEach((form) => {
            form.addEventListener(
              "submit",
              (event) => {
                if (!form.checkValidity()) {
                  event.preventDefault();
                  event.stopPropagation();
                }

                form.classList.add("was-validated");
              },
              false
            );
          });
        }
        formValidation();

        
      

      

      // Custom Alert
      function notify(msg, msgType) {
        notie.alert({
          type: msgType,
          text: msg,
        });
      }

       


      // Modal
      function notifyModal(tlt, txt, icn, cbt) {
        Swal.fire({
          title: tlt,
          text: txt,
          icon: icn,
          confirmButtonText: cbt,
        });
      }

      
      // Prompt
      function Prompt() {
        let toast = function (c) {
          const { msg = "", icon = "success", position = "top-end" } = c;
          const Toast = Swal.mixin({
            toast: true,
            position: position,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
              toast.addEventListener("mouseenter", Swal.stopTimer);
              toast.addEventListener("mouseleave", Swal.resumeTimer);
            },
          });

          Toast.fire({
            icon: icon,
            title: msg,
          });
        };
        let success = function (c) {
          const { msg = "", icon = "success", title = "", footer = "" } = c;
          Swal.fire({
            icon: icon,
            title: title,
            text: msg,
            footer: footer,
          });
        };

        let error = function (c) {
          const { msg = "", icon = "error", title = "", footer = "" } = c;
          Swal.fire({
            icon: icon,
            title: title,
            text: msg,
            footer: footer,
          });
        };

        async function custom(c) {
          const { msg = "",
                 title = "", 
                 icon="", 
                 showConfirmButton=true,
                 showCancelButton=true
                  } = c;

          const { value: result } = await Swal.fire({
            icon:icon,
            title: title,
            html: msg,
            backdrop: false,
            showConfirmButton:showConfirmButton,
            focusConfirm: false,
            showCancelButton: showCancelButton,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                                c.willOpen();
                            }
            },
            didOpen: () => {
              if (c.didOpen !== undefined) {
                                c.didOpen();
                            }
            },
            preConfirm: () => {
              if (c.preConfirm !== undefined) {
                                c.preConfirm();
                            }
            },
          });

          if (result) {
                    if (result.dismiss !== Swal.DismissReason.cancel) {
                        if (result.value !== "") {
                            if (c.callback !== undefined) {
                                c.callback(result);
                            }
                        } else {
                            c.callback(false);
                        }
                    } else {
                        c.callback(false);
                    }
                }
        }

        return {
          toast: toast,
          success: success,
          error: error,
          custom: custom,
        };
      }