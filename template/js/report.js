document.getElementById("fileInput").addEventListener("change", handleFileUpload);

function handleFileUpload(event) {
    const file = event.target.files[0]; // Получаем файл из input
    if (!file) {
        alert("Выберите файл для загрузки.");
        return;
    }

    const reader = new FileReader();
    reader.onload = async (e) => {
        const xmlString = e.target.result; // Содержимое файла
        const hosts = await parseNmapXML(xmlString); // Парсим XML
        renderHostsTable(hosts); // Отображаем таблицу
        renderPortDistributionChart(hosts); // Отображаем диаграмму
    };
    reader.readAsText(file); // Читаем файл как текст
}

async function parseNmapXML(xmlString) {
    const parser = new DOMParser();
    const xmlDoc = parser.parseFromString(xmlString, "application/xml");

    const hosts = [];
    const hostNodes = xmlDoc.getElementsByTagName("host");

    for (let host of hostNodes) {
        const ip = host.getElementsByTagName("address")[0]?.getAttribute("addr") || "N/A";
        const ports = host.getElementsByTagName("port");

        const openPorts = [];
        for (let port of ports) {
            const portNumber = port.getAttribute("portid");
            const state = port.getElementsByTagName("state")[0]?.getAttribute("state");
            const service = port.getElementsByTagName("service")[0]?.getAttribute("name") || "Unknown";
            if (state === "open") {
                openPorts.push({ port: portNumber, service });
            }
        }

        hosts.push({
            ip,
            openPorts,
        });
    }

    return hosts;
}

function renderHostsTable(hosts) {
    const tableContainer = document.getElementById("tableContainer");
    let html = `
        <table class="table table-striped">
            <thead>
                <tr>
                    <th>IP</th>
                    <th>Открытые порты</th>
                    <th>Службы</th>
                </tr>
            </thead>
            <tbody>
    `;

    hosts.forEach((host) => {
        const ports = host.openPorts.map((port) => port.port).join(", ");
        const services = host.openPorts.map((port) => port.service).join(", ");
        html += `
            <tr>
                <td>${host.ip}</td>
                <td>${ports}</td>
                <td>${services}</td>
            </tr>
        `;
    });

    html += "</tbody></table>";
    tableContainer.innerHTML = html;
}

function renderPortDistributionChart(hosts) {
    const portCounts = {};
    hosts.forEach((host) => {
        host.openPorts.forEach((port) => {
            portCounts[port.service] = (portCounts[port.service] || 0) + 1;
        });
    });

    const labels = Object.keys(portCounts);
    const data = Object.values(portCounts);

    const ctx = document.getElementById("portChart").getContext("2d");
    new Chart(ctx, {
        type: "pie",
        data: {
            labels: labels,
            datasets: [
                {
                    label: "Распределение портов",
                    data: data,
                    backgroundColor: ["#FF6384", "#36A2EB", "#FFCE56", "#4BC0C0"],
                },
            ],
        },
    });
}
