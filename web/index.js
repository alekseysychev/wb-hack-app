const SearchInput = document.getElementById('search')
const list =  document.getElementById('list')

function sendRequest(input) {
    console.log(input)
    return ['1', '2', '3']
}

SearchInput.onchange = () => {
    console.log('here')
    result = sendRequest(SearchInput.value)
    for (let i = 0; i < result.length; i++) {
        resultElement = result[i]
        const liRod = document.createElement('li')
        liRod.innerHTML = resultElement
        liRod.onclick = () => {
            SearchInput.value = liRod.innerHTML
        }
        list.appendChild(liRod)
    }
}
