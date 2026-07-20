
# 📄 Technical Specification: `UserModel`

> **Package:** models
> **Automatically generated documentation** by the Geanky tool.

---

## 1. Quick Summary (API & State)
A high-level overview of the class, its internal state, and available methods.

**Internal State & Dependencies:**


- `private ` **nome** (`String`)


- `private ` **idade** (`int`)


**Available Methods:**
- **setNome(String nome)** ➞ returns `void`
- **setIdade(int idade)** ➞ returns `void`
- **getIdade()** ➞ returns `int`
- **getNome()** ➞ returns `String`


---

## 2. Class Dependencies & State
Visual representation of the internal state and external dependencies this class maintains.

```mermaid
flowchart LR
    %% Styling
    classDef classNode fill:#2b3137,stroke:#fff,stroke-width:2px,color:#fff;
    classDef stateNode fill:#f4f6f8,stroke:#d0d7de,color:#24292f;
    classDef extNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    
    ThisClass["UserModel"]:::classNode

    %% State vs External Dependencies
    
    
    ThisClass -- "Maintains State" --- State_nome(["String<br>nome"]):::stateNode
    
    
    
    ThisClass -- "Maintains State" --- State_idade(["int<br>idade"]):::stateNode
    
    
```

---

## 3. Deep Dive (Constructors & Methods)
Expand the sections below to read the exact pseudo-code and business rules.


### 🛠️ Constructors

<details>
<summary><b>UserModel</b>(<i>String</i> nome, <i>int</i> idade) (Click to expand)</summary>

> **Signature:**
> `public UserModel(String nome, int idade)`

**Parameters:**

- **nome** (`String`)

- **idade** (`int`)


**Step-by-Step Logic:**



1. Set 'this.nome' to 'nome'

1. Set 'this.idade' to 'idade'



</details>




### ⚙️ Methods

<details>
<summary><b>setNome</b>(<i>String</i> nome) ➞ `void` (Click to expand)</summary>

> **Signature:**
> `public void setNome(String nome)`

**Data Flow:**
```mermaid
flowchart LR
    classDef methodNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    Caller(("Caller"))
    Method["setNome(String nome)"]:::methodNode

    Caller -- "Calls" --> Method
    Method -. "Returns<br>void" .-> Caller
```

**Step-by-Step Logic:**



1. Set 'this.nome' to 'nome'



**Parameters:**

- **nome** (`String`)


**Step-by-Step Logic:**



1. Set 'this.nome' to 'nome'



</details>

<details>
<summary><b>setIdade</b>(<i>int</i> idade) ➞ `void` (Click to expand)</summary>

> **Signature:**
> `public void setIdade(int idade)`

**Data Flow:**
```mermaid
flowchart LR
    classDef methodNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    Caller(("Caller"))
    Method["setIdade(int idade)"]:::methodNode

    Caller -- "Calls" --> Method
    Method -. "Returns<br>void" .-> Caller
```

**Step-by-Step Logic:**



1. Set 'this.idade' to 'idade'



**Parameters:**

- **idade** (`int`)


**Step-by-Step Logic:**



1. Set 'this.idade' to 'idade'



</details>

<details>
<summary><b>getIdade</b>() ➞ `int` (Click to expand)</summary>

> **Signature:**
> `public int getIdade()`

**Data Flow:**
```mermaid
flowchart LR
    classDef methodNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    Caller(("Caller"))
    Method["getIdade()"]:::methodNode

    Caller -- "Calls" --> Method
    Method -. "Returns<br>int" .-> Caller
```

**Step-by-Step Logic:**



1. Return the result of: this.idade



**Parameters:**
> *None.*


**Step-by-Step Logic:**



1. Return the result of: this.idade



</details>

<details>
<summary><b>getNome</b>() ➞ `String` (Click to expand)</summary>

> **Signature:**
> `public String getNome()`

**Data Flow:**
```mermaid
flowchart LR
    classDef methodNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    Caller(("Caller"))
    Method["getNome()"]:::methodNode

    Caller -- "Calls" --> Method
    Method -. "Returns<br>String" .-> Caller
```

**Step-by-Step Logic:**



1. Return the result of: this.nome



**Parameters:**
> *None.*


**Step-by-Step Logic:**



1. Return the result of: this.nome



</details>


