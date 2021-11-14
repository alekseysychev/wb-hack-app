const SearchInput = document.getElementById('search')
const list = document.getElementById('list')

function sendRequest(input) {
    var xmlhttp = new XMLHttpRequest();
    xmlhttp.onreadystatechange = function () {
        if (xmlhttp.readyState == XMLHttpRequest.DONE) {   // XMLHttpRequest.DONE == 4
            if (xmlhttp.status == 200) {
                let result = JSON.parse(xmlhttp.response);
                list.innerHTML = '';
                for (let i = 0; i < result.length; i++) {
                    const liRod = document.createElement('li')
                    liRod.innerHTML = result[i]
                    liRod.onclick = () => {
                        SearchInput.value = liRod.innerHTML;
                        sendRequest(SearchInput.value);
                    }
                    list.appendChild(liRod)
                }
            }
            else if (xmlhttp.status == 400) {
                alert('There was an error 400');
            }
            else {
                alert('something else other than 200 was returned');
            }
        }
    };
    xmlhttp.open("GET", "/api/query/" + SearchInput.value, true);
    xmlhttp.send();
}

SearchInput.onkeyup = () => {
    sendRequest(SearchInput.value);
}
