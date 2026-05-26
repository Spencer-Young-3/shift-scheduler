function calculateTotals(){
    day_elems = document.getElementsByClassName('day-row')
    week_total = 0
    Array.from(day_elems).forEach(element => {
        selected_elems = element.getElementsByClassName('selected')
        total = (selected_elems.length / 6).toFixed(2)
        week_total += Number(total)
        document.getElementById('day-total-' + element.dataset.day).innerHTML = total
    })
    document.getElementById('week-total-value').innerHTML = week_total
}

window.addEventListener("load", (event) => {
    calculateTotals()
})


/*
const form = document.querySelector("#schedule-form")
console.log("Attaching listeners to form")

form.addEventListener("htmx:beforeRequest", e => console.log("beforeRequest", e.detail))
form.addEventListener("htmx:beforeSwap", e => console.log("beforeSwap", e.detail, "shouldSwap:", e.detail.shouldSwap))
form.addEventListener("htmx:afterSwap", e => console.log("afterSwap", e.detail))
form.addEventListener("htmx:afterRequest", e => console.log("afterRequest", e.detail))
form.addEventListener("htmx:responseError", e => console.log("responseError", e.detail))
form.addEventListener("htmx:swapError", e => console.log("swapError", e.detail))
form.addEventListener("htmx:targetError", e => console.log("targetError", e.detail))

*/