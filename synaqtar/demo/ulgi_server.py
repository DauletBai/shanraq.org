#!/usr/bin/env python3
"""
Shanraq Template Server
Сервер для обработки шаблонов с компонентами
"""

import http.server
import socketserver
import os
import re
from urllib.parse import urlparse, unquote

class TemplateHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        # Parse the URL
        parsed_path = urlparse(self.path)
        path = unquote(parsed_path.path)
        
        # Handle root path
        if path == '/':
            path = '/index.html'
        
        print(f"Requested path: {path}")
        
        # Check if it's an HTML file
        if path.endswith('.html'):
            # Try to process ulgi
            ulgi_path = '.' + path
            print(f"Template path: {ulgi_path}")
            
            if os.path.exists(ulgi_path):
                try:
                    # Read ulgi content
                    with open(ulgi_path, 'r', encoding='utf-8') as f:
                        content = f.read()
                    
                    print(f"Original content length: {len(content)}")
                    
                    # Process components
                    content = self.process_components(content)
                    
                    print(f"Processed content length: {len(content)}")
                    
                    # Send response
                    self.send_response(200)
                    self.send_header('Content-type', 'text/html; charset=utf-8')
                    self.end_headers()
                    self.wfile.write(content.encode('utf-8'))
                    return
                except Exception as e:
                    print(f"Error processing ulgi {ulgi_path}: {e}")
        
        # Fall back to default behavior
        super().do_GET()
    
    def process_components(self, content):
        """Process {{> component_name }} syntax"""
        # Find all component references
        pattern = r'\{\{\s*>\s*(\w+)\s*\}\}'
        
        def replace_component(match):
            component_name = match.group(1)
            component_path = f'betjagy/bolıkter/{component_name}.html'
            
            print(f"Processing component: {component_name} from {component_path}")
            
            if os.path.exists(component_path):
                try:
                    with open(component_path, 'r', encoding='utf-8') as f:
                        component_content = f.read()
                        print(f"Loaded component {component_name}: {len(component_content)} chars")
                        return component_content
                except Exception as e:
                    print(f"Error loading component {component_name}: {e}")
                    return match.group(0)  # Return original if error
            else:
                print(f"Component {component_name} not found at {component_path}")
                return match.group(0)  # Return original if not found
        
        # Replace all component references
        content = re.sub(pattern, replace_component, content)
        return content

def run_server(port=8082):
    """Run the ulgi server"""
    handler = TemplateHandler
    
    with socketserver.TCPServer(("", port), handler) as httpd:
        print(f"🚀 Shanraq Template Server running on http://localhost:{port}")
        print(f"📄 Главная: http://localhost:{port}/")
        print(f"📝 Детальная: http://localhost:{port}/betjagy/better/index.html")
        print(f"🎯 Демо: http://localhost:{port}/betjagy/better/demo.html")
        print("=" * 50)
        print("Сервер остановить: Ctrl+C")
        
        try:
            httpd.serve_forever()
        except KeyboardInterrupt:
            print("\n🛑 Сервер остановлен")

if __name__ == "__main__":
    run_server()
