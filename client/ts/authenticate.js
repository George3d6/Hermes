const enableAuthenticationForm = () => {
    document.getElementById("submit_auth").addEventListener("click",(e) => {
        e.preventDefault();
        const identifier = document.getElementById("identifier_field").value;
        const credentials = document.getElementById("credentials_field").value;
        const xhr = new XMLHttpRequest();
        xhr.open("GET", `/authenticate/?identifier=${identifier}&credentials=${credentials}`);
        xhr.send();
        xhr.onreadystatechange = function() {
            $('#sign_in_form_modal')
              .modal('hide')
            ;
        }
    });
}

export default enableAuthenticationForm;
