
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

## 2. Architecture & Data Flow Diagram
Visual representation of how data enters the class, internal state, and external dependencies.

```mermaid
flowchart LR
    %% Styling
    classDef classNode fill:#2b3137,stroke:#fff,stroke-width:2px,color:#fff;
    classDef stateNode fill:#f4f6f8,stroke:#d0d7de,color:#24292f;
    classDef extNode fill:#0366d6,stroke:#fff,stroke-width:2px,color:#fff;
    
    Caller(("Caller"))
    ThisClass["UserModel"]:::classNode

    %% Method Calls
    
    Caller -- "Calls setNome(String nome)" --> ThisClass
    ThisClass -. "Returns void" .-> Caller
    
    Caller -- "Calls setIdade(int idade)" --> ThisClass
    ThisClass -. "Returns void" .-> Caller
    
    Caller -- "Calls getIdade()" --> ThisClass
    ThisClass -. "Returns int" .-> Caller
    
    Caller -- "Calls getNome()" --> ThisClass
    ThisClass -. "Returns String" .-> Caller
    

    %% State vs External Dependencies
    
    
    ThisClass -- "Maintains State" --- State_nome(["String nome"]):::stateNode
    
    
    
    ThisClass -- "Maintains State" --- State_idade(["int idade"]):::stateNode
    
    
```

---

## 3. Deep Dive (Constructors & Methods)
Expand the sections below to read the exact pseudo-code and business rules.


### 🛠️ Constructors

<details>
<summary><b>UserModel</b>(<i>String</i> nome, <i>int</i> idade) (Click to expand)</summary>

**Parameters:**

- **nome** (`String`)

- **idade** (`int`)


**Step-by-Step Logic:**



1. Set &#39;this.nome&#39; to &#39;nome&#39;

1. Set &#39;this.idade&#39; to &#39;idade&#39;



</details>




### ⚙️ Methods

<details>
<summary><b>setNome</b>(<i>String</i> nome) ➞ `void` (Click to expand)</summary>

> **Signature:** `public void setNome(String nome)`

**Parameters:**

- **nome** (`String`)


**Step-by-Step Logic:**



1. Set &#39;this.nome&#39; to &#39;nome&#39;



</details>

<details>
<summary><b>setIdade</b>(<i>int</i> idade) ➞ `void` (Click to expand)</summary>

> **Signature:** `public void setIdade(int idade)`

**Parameters:**

- **idade** (`int`)


**Step-by-Step Logic:**



1. Set &#39;this.idade&#39; to &#39;idade&#39;



</details>

<details>
<summary><b>getIdade</b>() ➞ `int` (Click to expand)</summary>

> **Signature:** `public int getIdade()`

**Parameters:**
> *None.*


**Step-by-Step Logic:**



1. Return the result of: this.idade



</details>

<details>
<summary><b>getNome</b>() ➞ `String` (Click to expand)</summary>

> **Signature:** `public String getNome()`

**Parameters:**
> *None.*


**Step-by-Step Logic:**



1. Return the result of: this.nome



</details>


