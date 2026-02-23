# -*- coding: utf-8 -*-
"""Location: ./plugins/deny_filter/deny.py
This is the public entry point that switches between Rust and Python.
"""

from mcpgateway.services.logging_service import LoggingService

# Setup logging
logging_service = LoggingService()
logger = logging_service.get_logger(__name__)

try:
    # 1. Attempt to load the Rust implementation from the neighbor file
    from .deny_rust import DenyListPluginRust as _impl
    RUST_AVAILABLE = True
    logger.info("DenyListPlugin: Rust implementation selected.")
except (ImportError, AttributeError):
    # 2. Fallback to the original Python file
    from .deny_orig import DenyListPlugin as _impl
    RUST_AVAILABLE = False
    logger.debug("DenyListPlugin: Rust not found, falling back to deny_orig.py")

DenyListPlugin = _impl

__all__ = ["DenyListPlugin", "RUST_AVAILABLE"]