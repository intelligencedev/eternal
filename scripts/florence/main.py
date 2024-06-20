# This script requires flash attention
# Example curl command:
# curl -X POST -H "Content-Type: multipart/form-data" -F "file=@/Users/art/Desktop/dollface.png" http://192.168.0.148:8081/generate
# Example response:
# {"<OD>":{"bboxes":[[0.25600001215934753,0.25600001215934753,511.2320251464844,511.2320251464844],[113.92000579833984,68.86400604248047,415.4880065917969,435.968017578125]],"labels":["doll","human face"]}}

from flask import Flask, request, jsonify
import requests
from PIL import Image
from transformers import AutoProcessor, AutoModelForCausalLM
import io

app = Flask(__name__)

# Load the model and processor
model = AutoModelForCausalLM.from_pretrained("microsoft/Florence-2-large-ft", trust_remote_code=True)
processor = AutoProcessor.from_pretrained("microsoft/Florence-2-large-ft", trust_remote_code=True)

@app.route('/generate', methods=['POST'])
def generate():
    # Get the image from the request
    if 'file' not in request.files:
        return jsonify({"error": "No file part in the request"}), 400

    file = request.files['file']
    if file.filename == '':
        return jsonify({"error": "No selected file"}), 400

    try:
        image = Image.open(file.stream).convert("RGB")  # Ensure image is in RGB format
    except Exception as e:
        return jsonify({"error": str(e)}), 400

    # Prepare the inputs
    prompt = "<OD>"
    inputs = processor(text=prompt, images=image, return_tensors="pt")

    # Generate the output
    generated_ids = model.generate(
        input_ids=inputs["input_ids"],
        pixel_values=inputs["pixel_values"],
        max_new_tokens=1024,
        do_sample=False,
        num_beams=3
    )
    generated_text = processor.batch_decode(generated_ids, skip_special_tokens=False)[0]

    # Post-process the generated text
    parsed_answer = processor.post_process_generation(generated_text, task="<OD>", image_size=(image.width, image.height))

    return jsonify(parsed_answer)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8081)
