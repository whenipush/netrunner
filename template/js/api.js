const API_BASE_URL = "http://localhost:3000/api/v1";

// Обработка хостов
const HostAPI = {
    async getAllHosts() {
        const response = await fetch(`${API_BASE_URL}/host`);
        if (!response.ok) throw new Error("Ошибка получения хостов");
        return response.json();
    },

    async createHost(host) {
        const response = await fetch(`${API_BASE_URL}/host`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(host),
        });
        if (!response.ok) throw new Error("Ошибка создания хоста");
        return response.json();
    },

    async updateHost(id, host) {
        const response = await fetch(`${API_BASE_URL}/host/${id}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(host),
        });
        if (!response.ok) throw new Error("Ошибка обновления хоста");
        return response.json();
    },

    async deleteHost(id) {
        const response = await fetch(`${API_BASE_URL}/host/${id}`, {
            method: "DELETE",
        });
        if (!response.ok) throw new Error("Ошибка удаления хоста");
    },
};

// Обработка групп
const GroupAPI = {
    async getAllGroups() {
        const response = await fetch(`${API_BASE_URL}/group`);
        if (!response.ok) throw new Error("Ошибка получения групп");
        return response.json();
    },

    async createGroup(group) {
        const response = await fetch(`${API_BASE_URL}/group`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(group),
        });
        if (!response.ok) throw new Error("Ошибка создания группы");
        return response.json();
    },

    async updateGroup(id, group) {
        const response = await fetch(`${API_BASE_URL}/group/${id}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(group),
        });
        if (!response.ok) throw new Error("Ошибка обновления группы");
        return response.json();
    },

    async deleteGroup(id) {
        const response = await fetch(`${API_BASE_URL}/group/${id}`, {
            method: "DELETE",
        });
        if (!response.ok) throw new Error("Ошибка удаления группы");
    },
};

// Обработка связи хостов и групп
const HostGroupAPI = {
    async addHostToGroup(data) {
        const response = await fetch(`${API_BASE_URL}/host-add-group`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
        });
        if (!response.ok) throw new Error("Ошибка добавления хоста в группу");
        return response.json();
    },
};

// Обработка Nmap
const NmapAPI = {
    async uploadScript(scriptData) {
        const response = await fetch(`${API_BASE_URL}/upload-script`, {
            method: "POST",
            body: scriptData, // Используем FormData
        });
        if (!response.ok) throw new Error("Ошибка загрузки скрипта");
        return response.json();
    },

    async processNmapRequest(data) {
        const response = await fetch(`${API_BASE_URL}/nmap`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
        });
        if (!response.ok) throw new Error("Ошибка выполнения Nmap");
        return response.json();
    },

    async getLastNmap() {
        const response = await fetch(`${API_BASE_URL}/last-nmap`);
        if (!response.ok) throw new Error("Ошибка получения последнего результата Nmap");
        return response.json();
    },

    async getAllNmapResults() {
        const response = await fetch(`${API_BASE_URL}/all-nmap`);
        if (!response.ok) throw new Error("Ошибка получения всех результатов Nmap");
        return response.json();
    },
    async getAllTasks() {
        const response = await fetch(`${API_BASE_URL}/task-all`);
        if (!response.ok) throw new Error("Ошибка получения всех результатов Nmap");
        return response.json();
    },
};

export { HostAPI, GroupAPI, HostGroupAPI, NmapAPI };
