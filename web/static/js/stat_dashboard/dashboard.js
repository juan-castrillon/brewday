const ctx = document.getElementById('myChart');
const data = window.Stats;
new Chart(ctx, {
    type: 'bar',
    data: {
        labels: data.map(row => row.FinishedTimeString),
        datasets: [
            {
                label: 'Evaporation',
                data: data.map(row => row.Evaporation)
            }
        ]
    }
});