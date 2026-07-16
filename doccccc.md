# 📦 <code>public</code> UsuarioController
 
> Documentação gerada automaticamente a partir da análise estática do código-fonte.
 
<p>
<img src="https://img.shields.io/badge/-public-6f42c1?style=flat-square" alt="public" /> <img src="https://img.shields.io/badge/fields-2-informational?style=flat-square" /> <img src="https://img.shields.io/badge/constructors-1-informational?style=flat-square" /> <img src="https://img.shields.io/badge/methods-1-informational?style=flat-square" />
</p>
 
---
 
## 🧩 Atributos
 
| Modificadores | Tipo | Nome |
|:---|:---|:---|
|  | <code>String</code> | <code>nomeSistema</code> |
|  | <code>int</code> | <code>ano</code> |

 
---
 
## 🏗️ Construtores
 
### public UsuarioController(String nomeSistema, int ano)
 
**📥 Contrato**
 
| Parâmetro | Tipo | Modificadores |
|:---|:---|:---|
| <code>nomeSistema</code> | <code>String</code> |  |
| <code>ano</code> | <code>int</code> |  |

**🚫 Retorna:** nada (<code>void</code>)
 
**⚙️ Comportamento**
 
- [x] this.nomeSistema = nomeSistema
- [x] this.ano = ano

<br>
 
---
 


## ⚙️ Métodos
 
### public processarUsuario(int idade, String status)
 
**📥 Contrato**
 
| Parâmetro | Tipo | Modificadores |
|:---|:---|:---|
| <code>idade</code> | <code>int</code> |  |
| <code>status</code> | <code>String</code> |  |

**↩️ Retorna:** <code>boolean</code>
 
**⚙️ Comportamento**
 
- [x] if (idade >= 18 && status.equals("ATIVO")) {
             System.out.println("Usuário processado no " + this.nomeSistema + "ano: " + ano);
             return true;
         }
- [x] return false

<br>
 
---
 


