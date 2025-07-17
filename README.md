# MegaETH

## 🤖 Modules
- Faucet
- Mint Start Fun NFT

## 📋 System Requirements

- Golang compiler
- [Capmonster API Key](https://capmonster.cloud/en)

## 🛠️ Installation

1. Install Golang (if it hasn't installed yet):

[Golang](https://go.dev/dl/)

2. Clone the repository
```bash
git clone [repository URL]
```

3. Install dependencies:
```bash
go mod tidy
```

## ⚙️ Configuration

### 1. Structure of configuration files

#### config/private_keys.txt
Private keys of EVM wallets (mandatory)

#### config/proxies.txt
The software supports different proxy formats. Don't forget to specify scheme - http:// https:// socks4:// socks5://

### settings.yaml Configuration

Edit the `config/config.yaml` file with the following settings:

```yaml
shuffle_accs: true # optional

delay_before_start:
  min: 10
  max: 30

delay_between_accs:
  min: 30
  max: 60

capmonster_api_key: "" # mandatory
```

4. Launch software:
```bash
go run main.go
```

## 🔒 Security Recommendations

1. **Protect Private Keys**: 
   - Never share your private keys or mnemonic phrases
   - Store sensitive data in secure, encrypted locations
   - Use environment variables or secure configuration management

2. **Proxy Usage**:
   - Use reliable and secure proxy servers
   - Rotate proxies to avoid IP blocking
   - Validate proxy credentials and connectivity

## 🤝 Contributing

### How to Contribute

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📜 License

Distributed under the MIT License. See `LICENSE` for more information.

## 📞 Support

For questions, issues, or support, please contact us through our Telegram channels.
