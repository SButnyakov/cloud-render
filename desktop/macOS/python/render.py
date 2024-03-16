import bpy

bpy.context.scene.render.image_settings.file_format = 'PNG'
bpy.context.scene.render.filepath = "E:/projects/GitHub/render-service/render-app/frame000.png"
bpy.ops.render.render(animation=False, write_still=True)
