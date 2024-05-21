package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Cliente representa os dados de um cliente
type Cliente struct {
	ID    string `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

var (
	clientes   = make(map[string]Cliente)
	clientesMu sync.Mutex
)

// Cria um novo cliente
func criarCliente(w http.ResponseWriter, r *http.Request) {
	var cliente Cliente
	err := json.NewDecoder(r.Body).Decode(&cliente)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clientesMu.Lock()
	clientes[cliente.ID] = cliente
	clientesMu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cliente)
}

// Obtém a lista de clientes
func obterClientes(w http.ResponseWriter, r *http.Request) {
	clientesMu.Lock()
	defer clientesMu.Unlock()

	var listaClientes []Cliente
	for _, cliente := range clientes {
		listaClientes = append(listaClientes, cliente)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listaClientes)
}

// Obtém um cliente por ID
func obterCliente(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	clientesMu.Lock()
	cliente, existe := clientes[id]
	clientesMu.Unlock()

	if !existe {
		http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cliente)
}

// Atualiza os dados de um cliente
func atualizarCliente(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var clienteAtualizado Cliente
	err := json.NewDecoder(r.Body).Decode(&clienteAtualizado)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clientesMu.Lock()
	cliente, existe := clientes[id]
	if !existe {
		clientesMu.Unlock()
		http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		return
	}

	cliente.Nome = clienteAtualizado.Nome
	cliente.Email = clienteAtualizado.Email
	clientes[id] = cliente
	clientesMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cliente)
}

// Deleta um cliente por ID
func deletarCliente(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	clientesMu.Lock()
	_, existe := clientes[id]
	if existe {
		delete(clientes, id)
	}
	clientesMu.Unlock()

	if !existe {
		http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/clientes", criarCliente).Methods("POST")
	r.HandleFunc("/clientes", obterClientes).Methods("GET")
	r.HandleFunc("/clientes/{id}", obterCliente).Methods("GET")
	r.HandleFunc("/clientes/{id}", atualizarCliente).Methods("PUT")
	r.HandleFunc("/clientes/{id}", deletarCliente).Methods("DELETE")

	http.Handle("/", r)
	fmt.Println("Servidor iniciado na porta 8080")
	http.ListenAndServe(":8080", nil)
}
