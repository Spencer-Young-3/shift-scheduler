function calculateTotals(){
    day_elems = document.getElementsByClassName('day-row')
    week_total = 0
    Array.from(day_elems).forEach(element => {
        selected_elems = element.getElementsByClassName('selected')
        total = (selected_elems.length / 6).toFixed(2)
        week_total += Number(total)
        document.getElementById('day-total-' + element.dataset.day).innerHTML = total
    })
    document.getElementById('week-total-value').innerHTML = week_total.toFixed(2)
}

window.addEventListener("load", (event) => {
    calculateTotals()
})
