from diffusers import StableDiffusionPipeline, DiffusionPipeline, DPMSolverMultistepScheduler
import torch

# Replace 'path/to/your/model.safetensors' with the actual path to your local .safetensors file
local_model_path = '/Users/arturoaquino/.eternal-v1/models/dreamshaper-8-turbo-sdxl/DreamShaperXL_Turbo_V2-SFW.safetensors'

pipe = StableDiffusionPipeline.from_single_file(local_model_path, torch_dtype=torch.float16, variant="fp16")
pipe.scheduler = DPMSolverMultistepScheduler.from_config(pipe.scheduler.config)
pipe = pipe.to("mps")

prompt = "portrait photo of muscular bearded guy in a worn mech suit, light bokeh, intricate, steel metal, elegant, sharp focus, soft lighting, vibrant colors"

generator = torch.manual_seed(0)
image = pipe(prompt, num_inference_steps=6, guidance_scale=2).images[0]  
image.save("./image.png")