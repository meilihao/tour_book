set -gx VULKAN_SDK /home/chen/opt/vulkan/x86_64
set -gx PATH $VULKAN_SDK/bin $PATH
set -gx LD_LIBRARY_PATH $VULKAN_SDK/lib $LD_LIBRARY_PATH
set -gx VK_LAYER_PATH $VULKAN_SDK/etc/explicit_layer.d
