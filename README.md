# **Hashcat Automation Tool**

This tool automates **Hashcat-based password cracking** using a structured workflow. It processes hash lists, wordlists, rules, and additional wordlists while ensuring all temporary and output files are stored in a **cache directory** (`cache/`).

---

## **Features**
✔ Automates **Hashcat cracking** with multiple attack modes  
✔ Supports **CeWL** for website-based wordlist generation  
✔ Handles **custom potfiles** and extracted passwords  
✔ Uses **rules-based cracking** for enhanced password recovery  
✔ Processes **large wordlists** with on-the-fly decompression (`7z`, `bzip2`)  
✔ Stores all outputs in `cache/` for easy access  

---

## **Configuration File: `config/config.json`**
The **`config.json`** file defines default paths for Hashcat, wordlists, rules, and cache storage. Modify it as needed before running the tool.

📌 **Ensure the `cache/` directory exists before running the tool.**  
📌 **Use `config.json.sample` as a template for your configuration.**  

---

## **Custom Wordlists, Rules, and Resources**
Use this section to list your **custom wordlists, rules, and additional resources**, including **URLs** for downloading wordlists.

### **Wordlists**
- 🔗 [Rockyou Wordlist](https://weakpass.com/wordlists/rockyou.txt)
- 🔗 [Passphrases Wordlist](https://github.com/initstring/passphrase-wordlist/releases/download/v2022.1/passphrases.txt)
- 🔗 [SecLists Passwords](https://github.com/danielmiessler/SecLists/tree/master/Passwords)
- 🔗 [CrackStation Wordlist](https://crackstation.net/buy-crackstation-wordlist-password-cracking-dictionary.htm)

### **Rules**
- 🔗 [clem9669_large.rule](https://github.com/clem9669/hashcat-rule/blob/master/clem9669_large.rule)
- 🔗 [rules_full.rule](https://github.com/Unic0rn28/hashcat-rules/blob/main/rules_full.7z)
- 🔗 [Passphrases Rules](https://github.com/initstring/passphrase-wordlist/tree/master/hashcat-rules)
- 🔗 [Hashcat Rules Collection](https://github.com/hashcat/hashcat/tree/master/rules)

### **Additional Resources**
- 🔗 [HashMob Wordlists](https://hashmob.net/resources/hashmob)
- 🔗 [WeakPass Wordlists](https://weakpass.com/wordlists)
- 🔗 [Kali Linux Wordlists](https://gitlab.com/kalilinux/packages/wordlists)

---

## **How to Build and Run the Tool**

### **1️⃣ Build the Binary**
To build the tool, navigate to the project directory and run:
```sh
go build -o hashcat-auto main.go
```

### **2️⃣ Run the Tool**
Execute with required parameters:
```sh
./hashcat-auto --hashlist=myhashes.txt --mode=0
```

### **3️⃣ Custom Example with CeWL and Additional Wordlists**
```sh
./hashcat-auto --hashlist=myhashes.txt --mode=1000 --wordlist=mywordlist.txt --url=https://example.com --enable-additional-wordlists
```

---

## **License**
📜 MIT License – Feel free to modify and improve the tool! 🚀

---

## **Contributors**
👨‍💻 **Maintainer:** Chris Meistre  
💬 **Issues?** Open an issue in the repository!
