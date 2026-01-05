# Deployment Guide (Ubuntu 24 LTS)

Your application is ready! Follow these steps to deploy to your VPS.

## 1. Upload the Binary
Upload the **`finance-app-linux`** binary to your VPS. You can use SCP, SFTP, or just drag-and-drop if your terminal supports it.

```bash
# Example using SCP (run from Windows)
scp .\finance-app-linux user@your-vps-ip:/home/user/
```

## 2. Setup on VPS
SSH into your VPS and set up the permissions.

```bash
# Make it executable
chmod +x finance-app-linux

# (Optional) Allow it to bind to port 80 if you want to run it directly without Nginx
# sudo setcap cap_net_bind_service=+ep ./finance-app-linux
```

## 3. Run It
You can run it directly to test:
```bash
./finance-app-linux
```
Your app will start on port **8081**.
Access it at: `http://your-vps-ip:8081`

---

## 4. (Recommended) Run in Background with SystemD
To keep the app running even after you close the terminal, create a service.

1. Create a service file:
   ```bash
   sudo nano /etc/systemd/system/finance.service
   ```

2. Paste this configuration (adjust paths if needed):
   ```ini
   [Unit]
   Description=Finance App
   After=network.target

   [Service]
   User=root
   # Change this directory to where you uploaded the file
   WorkingDirectory=/home/user
   # The command to start the app
   ExecStart=/home/user/finance-app-linux
   Restart=always

   [Install]
   WantedBy=multi-user.target
   ```

3. Start and Enable:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl start finance
   sudo systemctl enable finance
   ```

4. Check Log:
   ```bash
   journalctl -u finance -f
   ```

## 5. Security Notes
- The database `database.db` will be created automatically in the same folder.
- `assets` and `templates` are embedded inside the binary, so you **do not** need to upload them separately.
- **Default Admin**: `admin` / `admin123` (Change this immediately after logging in!)
